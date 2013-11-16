package main

import(
    "fmt"
    "regexp"
    "io/ioutil"
    "io"
    "net/http"
    "strings"
    iconv "github.com/djimenez/iconv-go"
    "os"
    "bufio"
)


var (
    url string
    src string
    err error
)

func main(){

    //逐行读取域名列表
    f , err := os.Open("dm.txt")
    if err != nil{
       fmt.Printf("%v\n",err)
       os.Exit(1)
    }

   defer f.Close()
   r := bufio.NewReader(f)
   
   //调用spider获取源码
   for {
        line, _, err := r.ReadLine()

        if err == io.EOF {
            fmt.Println("Spider end.")
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            break
        }
        src, err = spider("http://" + string(line))
        //fmt.Println(string(line)) 
        //fmt.Println(src)
   }


}


//页面爬虫
func spider(url string)(src string, err error) {

    //获取源码
    //url := "http://www.163.com/"

    resp, err := http.Get(url)

    if err != nil {
        fmt.Println("http get error.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("http read error")
    }

    src = string(body)


    //获取编码
    //src = `<html><head><meta charset="gb2312"></head><body>abc</body></html>`
    var myExp = myRegexp{regexp.MustCompile(`charset="(?P<c>[^\"]*)`)}
    re1 := myExp.FindStringSubmatchMap(src)
    charset := re1["c"]
    //fmt.Println(charset)

    //转换编码
    src,_ = iconv.ConvertString(src, charset, "utf-8")

    //将HTML标签全转换成小写
    re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
    src = re.ReplaceAllStringFunc(src, strings.ToLower)

    //去除STYLE
    re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
    src = re.ReplaceAllString(src, "")

    //去除SCRIPT
    re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
    src = re.ReplaceAllString(src, "")

    //去除所有尖括号内的HTML代码，并换成换行符
    re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
    src = re.ReplaceAllString(src, "\n")

    //去除连续的换行符
    re, _ = regexp.Compile("\\s{2,}")
    src = re.ReplaceAllString(src, "\n")

    //去掉空格
    re, _ = regexp.Compile("\\&nbsp;")
    src = re.ReplaceAllString(src, "")

    //去掉<>
    re, _ = regexp.Compile("\\&[l|g]t;")
    src = re.ReplaceAllString(src, "")

    //fmt.Println(src)
    return

}


//正则扩展，用于爬虫获取页面编码
type myRegexp struct{ 
    *regexp.Regexp
}
 
func(r *myRegexp)FindStringSubmatchMap(s string) map[string]string{
    captures:=make(map[string]string)
 
    match:=r.FindStringSubmatch(s)
    if match==nil{
        return captures
    }
 
    for i,name:=range r.SubexpNames(){
        if i==0||name==""{
            continue
        }
         captures[name]=match[i]
 
    }
    return captures
}