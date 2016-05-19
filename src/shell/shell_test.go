package shell


import (
    
    "testing"
    "text/template"
    "os"
    "reflect"
    "fmt"
    "strings"
    // "path"
    "bufio"
)

func TestReadFile(t *testing.T){
    // path.Join("../tmpl")
    // tmpl := template.Must( template.New("shell").Parse("/Users/tangshiyou/go_project/goshell/src/tmpl/file.tmpl"))
    tmpl := template.Must( template.ParseFiles("/Users/tangshiyou/go_project/goshell/src/tmpl/file.tmpl"))
    
    
    tmpl.Execute(os.Stdout, ShellScan{Pkgs:[]string{"fmt","os"}})
    fmt.Println(" data transfer :::: ")
    fmt.Println((&ShellScan{Pkgs:[]string{"fmt","os"}}).transferToTmplData())
}

func TestScan(t *testing.T){
    // ss := &ShellScan{}
    br := bufio.NewReader(os.Stdin)
    line , err := br.ReadSlice('\n')
    fmt.Println(line, err)
    // i := 0
    // for {
    //     i ++
    //     if i ==6 {
    //         break
    //     }
    //     fmt.Print(">>>")
    //     line, err := br.ReadString('\n')
    //     fmt.Println("line :: err ;; ",line,err)
        
    //     convertedStr := strings.Trim(line, " ")
    //     if len(convertedStr) == 1 && strings.HasPrefix(convertedStr,"$"){
    //         // fmt.Println("show key word ", showKeywords())
    //     }
        
    //     if strings.HasPrefix(strings.TrimLeft(line, " "), "$RUN"){
    //         // run 
    //         // fmt.Println(ss.transferToTmplData)
    //     }
        
    //     if strings.HasPrefix(strings.TrimLeft(line, " "), "$Q"){
    //         return 
            
    //     }
        
    //     // ss.MainBody = append(ss.MainBody, line)
    // }
    
}

