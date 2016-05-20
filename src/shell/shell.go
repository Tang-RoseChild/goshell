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

var (
    stepSetting bool
)

var (
    blocked bool // use for scanning a block statement,like type, func and so on
    blockedString string 
    blockAppend *[]string
)

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
        if blocked {
        fmt.Print(">>>>>>")
        } else {
            fmt.Print(">>>")
        }
        
        line, err := br.ReadString('\n')
        if err != nil {
            continue
        }
        
        if blocked{
            blockedString += line
            if strings.Contains(line,"}"){
                blocked = false
                
                *blockAppend = append(*blockAppend, blockedString)
                blockedString = ""
                blockAppend = nil
            } 
            
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
        if strings.HasPrefix(strings.TrimLeft(line, " "), "$CLEAR"){
            // print
            *ss = ShellScan{}
            
        }
        
        if strings.HasPrefix(strings.TrimLeft(line, " "), "$STEP"){
            
            split := strings.Split(line, " ")
            
            if len(split) == 2 {
                trimedStr := strings.Trim(split[1],"\n" )
                
                 if  trimedStr == "ON"{
                     
                     stepSetting = true
                     fmt.Println("ON :: ",stepSetting)
                     
                 } else if  trimedStr == "OFF"{
                     stepSetting = false
                 }
            }
            
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
        blocked = true
        blockAppend = &(ss.TypeBody)
        typeStr := strings.TrimRight(strings.TrimLeft(line, " ")," ")
        if strings.Contains(typeStr,"{") == false && strings.Contains(typeStr,"}")  == false {
            blocked = false
            blockAppend = nil
            ss.TypeBody = append(ss.TypeBody, line)
        } else {
            blockedString = line
        }
        
    case strings.HasPrefix(trimedStr, "func"):
        funcStr := strings.TrimRight(strings.TrimRight(strings.TrimLeft(line, " ")," "),"\n")
        if strings.HasSuffix(funcStr,"{") == false{
            fmt.Println("func statement should ends with {")
        } else {
            blocked = true
            blockedString += line
            blockAppend = &(ss.FuncBody)
        }
            
        
    case strings.HasPrefix(trimedStr, "global"):
        ss.TypeBody = append(ss.GlobalVars, line)
    case strings.HasPrefix(trimedStr, "import"):
        ss.TypeBody = append(ss.Pkgs, line)
    case strings.HasPrefix(trimedStr, "const"):
        ss.TypeBody = append(ss.ConstBody, line)
    
    default :
        ss.MainBody = append(ss.MainBody, line)    
    }
    
    if stepSetting {
        fmt.Println(" step true")
        ss.run()
    }
}

func (ss *ShellScan)run(){
    gf := ss.makeGoFile()
    defer func(){
        gf.Close()
        os.Remove(gf.Name())
    }()
    
    
    handleComOutErr( exec.Command("go", "run", gf.Name()).CombinedOutput() )
    
    
    
    
    
}
func (ss *ShellScan)makeGoFile() *os.File{
    tmpl := template.Must( template.ParseFiles("../tmpl/file.tmpl"))
    f, err := ioutil.TempFile(".", "temp")
    handleErr(err)
    defer f.Close()
    tmpl.Execute(f, ss.transferToTmplData())
    newPath := f.Name() + ".go"
    os.Rename(f.Name(), newPath)
    
    handleComOutErr( exec.Command("gofmt", "-w",newPath).CombinedOutput() )
    
    handleComOutErr( exec.Command("goimports", "-w",newPath).CombinedOutput())
    
    nf, err := os.Open(newPath)
    handleErr(err)
    return nf
}

func (ss *ShellScan)print(){
    gf := ss.makeGoFile()
    defer func(){
        gf.Close()
        os.Remove(gf.Name())
    }()
    fmt.Println(gf.Name())
    
    var count int 
    
    bf := bufio.NewReader(gf)
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
func handleComOutErr(bytes []byte, err error  ){
     if _, ok := err.(*exec.ExitError);!ok{
        handleErr(err)
    }
    fmt.Println(string(bytes))
}