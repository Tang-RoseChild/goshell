package shell 

import (
    // "testing"
    // "text/template"
    "os"
    "reflect"
    "fmt"
    "strings"
    // "path"
    "bufio"
    "text/template"
    "io/ioutil"
    "os/exec"
    "io"
)


const (
    contact = ";"
)

type usage string 
var keywords = map[string]usage{
    "$":usage("list all keyword"),
    "$Q":usage("quit"),
    "$RUN":usage("run all input"),
    "$STEP":usage("run step by step if flag is on"),
    "$CLEAR":usage("clear all input"),
}

var SS *ShellScan = &ShellScan{}
var TmplData *tmplData = &tmplData{}

type ShellScan struct{
    Pkgs []string
    TypeBody []string 
    GlobalVars []string 
    ConstBody []string 
    MainBody []string 
    FuncBody []string 
}

type tmplData struct{
    Pkgs string
    TypeBody string 
    GlobalVars string 
    ConstBody string 
    MainBody string 
    FuncBody string 
}

func (s *ShellScan) transferToTmplData () tmplData{
    data := &tmplData{}
    dataRv := reflect.ValueOf(data).Elem()
    rv := reflect.ValueOf(*s)
    rt := rv.Type()
    for i := 0 ; i < rv.NumField(); i ++ {
        f := rv.Field(i)
        slice := f.Interface().([]string)
        dataField := dataRv.FieldByName(rt.Field(i).Name)
        dataField.SetString(strings.Join(slice,contact))
        
    }
    return *data
}

func showKeywords()string {
    var ret []string
    for k, use := range keywords{
        ret = append(ret, k + "\t" ,string(use))
    }
    
    return strings.Join(ret, "\n")
}

func HandleScan(){
       br := bufio.NewReader(os.Stdin)
    
    
    for {
    
        fmt.Print(">>>")
        line, err := br.ReadString('\n')
        if err != nil {
            continue
        }
       
        if  strings.HasPrefix(strings.TrimLeft(line, " "), "$"){
        code := SS.operateHandle(line)
        if code == 1 {
            break
        }
        } else {
            SS.cmdFileHandle(line)
        }
        
        
    }
}

func (ss *ShellScan) operateHandle(line string) int{
     convertedStr := strings.Trim(line, " ")
    if len(convertedStr) == 2 && strings.HasPrefix(convertedStr,"$"){
            fmt.Println(showKeywords())
        }
        
        if strings.HasPrefix(strings.TrimLeft(line, " "), "$RUN"){
            // run 
            ss.run()
            ss.MainBody = nil
        }
        
        if strings.HasPrefix(strings.TrimLeft(line, " "), "$Q"){
            return 1
            
        }
        
        if strings.HasPrefix(strings.TrimLeft(line, " "), "$PRINT"){
            // print
            ss.print()
        }
        
        return 0
    }
func (ss *ShellScan) cmdFileHandle(line string) {
    // type 
    // func 
    // global vars 
    // const 
    // pkg 
    // main body 
    
    trimedStr := strings.TrimLeft(line, " ")
    
    switch  {
    case strings.HasPrefix(trimedStr, "type"):
        ss.TypeBody = append(ss.TypeBody, line)
    
    default :
        ss.MainBody = append(ss.MainBody, line)    
    }
}

func (ss *ShellScan)run(){
    tmpl := template.Must( template.ParseFiles("../tmpl/file.tmpl"))
    f, err := ioutil.TempFile(".", "temp")
    handleErr(err)
    defer f.Close()
    tmpl.Execute(f, ss.transferToTmplData())
    newPath := f.Name() + ".go"
    os.Rename(f.Name(), newPath)
    defer os.Remove(newPath)
    
    handleErr(exec.Command("gofmt", "-w",newPath).Run())
    handleErr(exec.Command("goimports", "-w",newPath).Run())
    out, err := exec.Command("go", "run", newPath).CombinedOutput()
    
    if _, ok := err.(*exec.ExitError);!ok{
        handleErr(err)
    }
    fmt.Println(string(out))
    
    
    
}

func (ss *ShellScan)print(){
    tmpl := template.Must( template.ParseFiles("../tmpl/file.tmpl"))
    f, err := ioutil.TempFile(".", "temp")
    handleErr(err)
    defer f.Close()
    tmpl.Execute(f, ss.transferToTmplData())
    newPath := f.Name() + ".go"
    os.Rename(f.Name(), newPath)
    defer os.Remove(newPath)
    
    handleErr(exec.Command("gofmt", "-w",newPath).Run())
    handleErr(exec.Command("goimports", "-w",newPath).Run())
    nf, err := os.Open(newPath)
    handleErr(err)
    defer nf.Close()
    var count int 
    
    bf := bufio.NewReader(nf)
    for {
        line, err := bf.ReadString('\n')
        if err == io.EOF {
            break
        }
        handleErr(err)
        count++
        fmt.Print(count," : ",line)
    }
    
    
    
}

func handleErr( err error){
    if err != nil {
        panic(err)
    }
}