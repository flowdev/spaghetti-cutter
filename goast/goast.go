package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	goTestFileName = `_test.go`
)

//
// packageDict is a simple dictionary of all known packages/paths and
// their source parts.
//

type goPackage struct {
	path       string
	deps       map[string]struct{} // or: map[string]*goPackage
	complexity int
}

type packageDict struct {
	root  string
	packs map[string]*goPackage
}

func newPackageDict(root string) *packageDict {
	return &packageDict{
		root:  root,
		packs: make(map[string]*goPackage),
	}
}

func (pd *packageDict) addPackage(path string, deps []string, complexity int) {
	pd.packs[path] = &goPackage{path: path, partMap: partMap}
}

func (pd *packageDict) addFileDeps(
	astImps []*ast.ImportSpec,
	packDict *packageDict,
	fset *token.FileSet,
) {
	imps := make(map[string]string)
	for _, astImp := range astImps {
		val := strings.Trim(astImp.Path.Value, "\"")
		pd.deps[val] = struct{}{}
	}
}

//
// Parse and process the main directory/package.
//

// WalkDirTree is walking the directory tree starting at the given root path
// looking for Go packages and analyzing them.
func WalkDirTree(root string) error {
	packDict := newPackageDict(root)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("ERROR: While walking the path %q: %v", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "vendor" {
			log.Printf("INFO: skipping vendor dir: %s", path)
			return filepath.SkipDir
		}
		if info.IsDir() {
			return processDir(path, packDict)
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR: While walking the root %q: %v", root, err)
		return
	}
}

func processDir(dir string, packDict *packageDict) error {
	fset := token.NewFileSet() // needed for any kind of parsing
	fmt.Println("Parsing the whole directory:", dir)
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("unable to parse the directory '%s': %v", dir, err)
	}
	for _, pkg := range pkgs { // iterate over subpackages (e.g.: xxx and xxx_test)
		if err := processPackage(pkg, fset, packDict); err != nil {
			return err
		}
	}
	return nil
}

// processPackage is processing all the files of one Go package.
func processPackage(pkg *ast.Package, fset *token.FileSet, packDict *packageDict) error {
	fmt.Println("processing package:", pkg.Name)
	partMap := make(map[string]*sourcePart)
	flows := make([]*sourcePart, 0, 128)
	fileMap := make(map[string]*mdFile)
	var err error

	for name, astf := range pkg.Files {
		fImps := newFileImps(astf.Imports, packDict, fset)
		baseName := goNameToBase(name)
		fileMap[baseName] = &mdFile{name: baseName, fImps: fImps}
		if flows, err = findSourceParts(
			partMap, flows,
			astf,
			name, "", fset,
		); err != nil {
			return fmt.Errorf(
				"unable to find all flows in package (%s): %v", pkg.Name, err)
		}
	}
	fmt.Println("Found", len(flows), "flows.")
	for _, f := range flows {
		if err = startFlowFile(f, fileMap); err != nil {
			return fmt.Errorf(
				"unable to start all Markdown files in package (%s): %v",
				pkg.Name, err)
		}
		if err = addToMDFile(f, partMap); err != nil {
			return fmt.Errorf(
				"unable to process all flows in package (%s): %v", pkg.Name, err)
		}
	}
	fmt.Println("processed flows with ", len(partMap), "souce parts.")
	for _, f := range fileMap {
		if err = endMDFile(f); err != nil {
			log.Printf("Error while ending file: %v", err)
		}
	}
	fmt.Println("Ended", len(fileMap), "files.")
	return nil
}

//
// Handle source file
//

func findSourceParts(
	partMap map[string]*sourcePart, flows []*sourcePart,
	astf *ast.File,
	goname string, path string, fset *token.FileSet,
) ([]*sourcePart, error) {
	baseName := goNameToBase(goname)

	for _, idecl := range astf.Decls {
		switch decl := idecl.(type) {
		case *ast.FuncDecl:
			doc := decl.Doc.Text()
			name := decl.Name.Name
			if strings.Contains(doc, flowMarker) {
				if i := strings.Index(name, "_"); i >= 0 {
					name = name[:i] // cut off the port name
				}
				flow := &sourcePart{
					kind:       sourcePartFlow,
					name:       name,
					doc:        doc,
					start:      lineFor(decl.Pos(), fset),
					end:        lineFor(decl.End(), fset),
					importPath: path,
					goFile:     goname,
					mdFile:     &mdFile{name: baseName},
				}
				partMap[markerFlow+name] = flow
				flows = append(flows, flow)
			} else {
				partMap[markerFunc+decl.Name.Name] = &sourcePart{
					kind:       sourcePartFunc,
					name:       name,
					start:      lineFor(decl.Pos(), fset),
					end:        lineFor(decl.End(), fset),
					importPath: path,
					goFile:     goname,
				}
			}
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				for _, s := range decl.Specs {
					ts := s.(*ast.TypeSpec)
					name := ts.Name.Name
					partMap[markerType+name] = &sourcePart{
						kind:       sourcePartType,
						name:       name,
						start:      lineFor(ts.Pos(), fset),
						end:        lineFor(ts.End(), fset),
						importPath: path,
						goFile:     goname,
					}
				}
			}
		}
	}
	return flows, nil
}
func goNameToBase(goname string) string {
	ext := filepath.Ext(goname)
	return goname[:len(goname)-len(ext)]
}
func lineFor(p token.Pos, fset *token.FileSet) int {
	if p.IsValid() {
		pos := fset.PositionFor(p, false)
		return pos.Line
	}

	return 0
}

//
// Write to Markdown file
//

func startFlowFile(flow *sourcePart, fileMap map[string]*mdFile) error {
	file := fileMap[flow.mdFile.name]
	if file == nil {
		return fmt.Errorf("missing flow file: " + flow.mdFile.name)
	}
	if file.osfile == nil {
		osfile, err := startMDFile(flow.mdFile.name)
		if err != nil {
			return err
		}
		file.osfile = osfile
	}
	flow.mdFile = file
	return nil
}

func startMDFile(fileBaseName string) (*os.File, error) {
	mdname := fileBaseName + ".md"

	f, err := os.Create(mdname)
	if err != nil {
		return nil, err
	}

	if _, err = f.WriteString(mdStart + fileBaseName + ".go\n\n"); err != nil {
		return nil, err
	}

	return f, nil
}

func addToMDFile(f *sourcePart, partMap map[string]*sourcePart) error {
	fmt.Println("processing flow:", f.name)
	if _, err := f.mdFile.osfile.WriteString(
		fmt.Sprintf(flowStart, f.name, f.goFile, f.start, f.end)); err != nil {

		return err
	}
	start, flow, end := ExtractFlowDSL(f.doc)
	if _, err := f.mdFile.osfile.WriteString(start + "\n"); err != nil {
		return err
	}
	log.Printf("Converting FlowDSL: '%s'\n", flow)
	svg, compTypes, dataTypes, info, err := gflowparser.ConvertFlowDSLToSVG(flow, f.name)
	if err != nil {
		return err
	}
	if info != "" {
		log.Printf("INFO: %s", info)
	}
	if err = ioutil.WriteFile(f.name+".svg", svg, os.FileMode(0666)); err != nil {
		return err
	}
	if _, err = f.mdFile.osfile.WriteString(
		fmt.Sprintf("![Flow: %s](./%s.svg)\n\n", f.name, f.name)); err != nil {

		return err
	}
	if err = writeReferences(f, compTypes, dataTypes, partMap); err != nil {
		return err
	}
	if _, err = f.mdFile.osfile.WriteString(end); err != nil {
		return err
	}

	return nil
}

func writeReferences(
	f *sourcePart, compTypes []data.Type,
	dataTypes []data.Type,
	partMap map[string]*sourcePart,
) error {
	dataTypes = filterTypes(dataTypes)
	dataTypes = sortTypes(dataTypes)
	compTypes = sortTypes(compTypes)
	dataLinks := getLinksForTypes(dataTypes, partMap, f.mdFile)

	n := max(len(compTypes), len(dataLinks))
	if n == 0 {
		return nil
	}

	if _, err := f.mdFile.osfile.WriteString(referenceTableHeader); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		row := bytes.Buffer{}
		if i < len(compTypes) {
			addComponentToRow(&row, compTypes[i], partMap, f.mdFile)
		}
		row.WriteString(" | ")
		if i < len(dataLinks) {
			addTypeToRow(&row, dataLinks[i])
		}
		row.WriteRune('\n')
		if _, err := f.mdFile.osfile.Write(row.Bytes()); err != nil {
			return err
		}
	}
	if _, err := f.mdFile.osfile.WriteString("\n"); err != nil {
		return err
	}
	return nil
}
func sortTypes(types []data.Type) []data.Type {
	sort.Slice(types, func(i, j int) bool {
		if types[i].Package == types[j].Package {
			return types[i].LocalType < types[j].LocalType
		}
		return types[i].Package < types[j].Package
	})
	return types
}
func filterTypes(types []data.Type) []data.Type {
	result := make([]data.Type, 0, len(types))
	for _, t := range types {
		if t.Package != "" {
			result = append(result, t)
			continue
		}
		s := t.LocalType
		if strings.HasPrefix(s, "[]") {
			s = s[2:]
		} else if strings.HasPrefix(s, "map[") {
			continue
		}
		switch s {
		case "bool", "byte", "complex64", "complex128", "float32", "float64",
			"int", "int8", "int16", "int32", "int64",
			"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64",
			"uintptr":
			continue
		default:
			t.LocalType = s
			result = append(result, t)
		}
	}
	return result
}
func getLinksForTypes(types []data.Type, partMap map[string]*sourcePart, mdFile *mdFile) []string {
	links := make([]string, 0, len(types))
	for _, typ := range types {
		link := getLinkForType(typ, partMap, mdFile)
		if link != "" {
			links = append(links, link)
		}
	}
	return links
}
func typeToString(t data.Type) string {
	if t.Package != "" {
		return t.Package + "." + t.LocalType
	}
	return t.LocalType
}
func addComponentToRow(row *bytes.Buffer, comp data.Type, partMap map[string]*sourcePart, mdFile *mdFile) {
	var flow, fun *sourcePart
	cNam := typeToString(comp)

	if comp.Package == "" {
		flow = partMap[markerFlow+cNam]
		fun = partMap[markerFunc+cNam]
	} else {
		flow = mdFile.fImps.getPartFor(comp.Package, markerFlow+comp.LocalType)
		fun = mdFile.fImps.getPartFor(comp.Package, markerFunc+comp.LocalType)
	}
	if flow != nil {
		fileName, err := fileNameFor(flow, markerFlow, mdFile)
		if err != nil {
			fmt.Println("WARNING: Unable to compute correct URL for flow", cNam, ":", err)
			fileName = flow.mdFile.name + ".md"
		}
		// [link to Google!](http://google.com)
		row.WriteString(
			"[" + cNam + "](" +
				fileName + "#flow-" +
				strings.ToLower(flow.name) +
				")")
	} else if fun != nil {
		fileName, err := fileNameFor(fun, markerFunc, mdFile)
		if err != nil {
			fmt.Println("WARNING: Unable to compute correct URL for function", cNam, ":", err)
			fileName = fun.goFile
		}
		row.WriteString(fmt.Sprintf(
			"[%s](%s#L%dL%d)",
			cNam, fileName, fun.start, fun.end,
		))
	} else {
		row.WriteString(cNam)
	}
}
func addTypeToRow(row *bytes.Buffer, tNam string) {
	row.WriteString(tNam)
}
func getLinkForType(typ data.Type, partMap map[string]*sourcePart, mdFile *mdFile) string {
	var ty *sourcePart
	tNam := typeToString(typ)
	if typ.Package == "" {
		ty = partMap[markerType+tNam]
	} else {
		ty = mdFile.fImps.getPartFor(typ.Package, markerType+typ.LocalType)
	}
	if ty == nil {
		return ""
	}

	fileName, err := fileNameFor(ty, markerType, mdFile)
	if err != nil {
		fmt.Println("WARNING: Unable to compute correct URL for type", tNam, ":", err)
		fileName = ty.goFile
	}
	return fmt.Sprintf(
		"[%s](%s#L%dL%d)",
		tNam, fileName, ty.start, ty.end,
	)
}
func fileNameFor(part *sourcePart, marker string, mdFile *mdFile) (string, error) {
	if marker == markerFlow {
		if mdFile.name == part.mdFile.name { // same MD file
			return "", nil
		}
		return outsideFileNameFor(part.mdFile.name+".md", part, mdFile)
	}

	return outsideFileNameFor(part.goFile, part, mdFile)
}
func outsideFileNameFor(name string, part *sourcePart, mdFile *mdFile) (string, error) {
	absF := name
	if !filepath.IsAbs(absF) {
		absF = filepath.Join(mdFile.fImps.packDict.cwd, name)
	}
	relF, err := filepath.Rel(mdFile.fImps.packDict.projRoot, absF)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(relF, ".."+string(filepath.Separator)) {
		return filepath.Rel(mdFile.fImps.packDict.cwd, absF) // inside of project always use relative paths
	}
	// outside of project:
	if mdFile.fImps.packDict.localLinks {
		return absF, nil
	}
	_, lastF := filepath.Split(absF)
	urlParts := strings.SplitN(part.importPath, "/", 4)
	url := "https://" + path.Join(urlParts[:3]...) + "/blob/master"
	if len(urlParts) > 3 {
		url += "/" + urlParts[3]
	}
	return url + "/" + lastF, nil
}

// ExtractFlowDSL extracts the flow DSL from a documentation comment string.
// The doc string should be given without comment characters.
// This function returns everything before the flow in start,
// the flow DSL itself and everything after it in end.
func ExtractFlowDSL(doc string) (start, flow, end string) {
	i := strings.Index(doc, flowMarker)
	if i < 0 {
		return doc, "", ""
	}
	start = doc[:i+1]
	i += len(flowMarker)

	buf := bytes.Buffer{}
	for dsl, ok := getDSLLine(doc, &i); ok; dsl, ok = getDSLLine(doc, &i) {
		buf.WriteString(dsl)
	}

	end = doc[i:]
	if end != "" && end[len(end)-1:] != "\n" {
		end += "\n"
	}
	return start, buf.String(), end
}
func getDSLLine(doc string, pi *int) (string, bool) {
	if *pi >= len(doc) {
		return "", false
	}
	tail := doc[*pi:]
	n := strings.IndexRune(tail, '\n')
	line := ""
	if n >= 0 {
		n++ // include NL
		line = tail[:n]
	} else {
		n = len(tail)
		line = tail + "\n" // add missing NL
	}

	dslN := len(dslMarker)
	if strings.TrimSpace(line) == "" { // support empty lines
		*pi += n
		return "\n", true
	}
	if n > dslN && line[:dslN] == dslMarker { // real DSL
		*pi += n
		return line[dslN:], true
	}
	return "", false
}

func endMDFile(f *mdFile) error {
	if f == nil || f.osfile == nil {
		return nil
	}
	return f.osfile.Close()
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
