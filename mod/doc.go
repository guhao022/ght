package mod

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//参数结构
type param struct {
	name string //参数名
	tp   string //参数类型
	call string //是否必须
	note string //注释
}

//生成文档内容结构
type docnote struct {
	tab         string //标记要生成文档的注释，包含在注释里
	title       string
	description string
	param       []*param
	success     string
	fail        string
	router      string
}

//标记文件内容结构
type signfile struct {
	version     string //版本
	title       string //文档名称
	description string //说明
	contact     string //联系人
	port        string //端口
	server      string //监听服务器
}

type Doc struct {
	root      string
	docpath   string //文档文件所存的位置
	extension string //扩展名
}

var signfilename = ".gdoc" //定义标记文件名称
var doctab = "__doc__"     //定义文档标记
var PthSep = string(os.PathSeparator)

//遍历文件夹 获取标记文件地址
func getSignPath(path string) []string {
	reg, err := regexp.Compile(signfilename)
	if err != nil {
		panic(err)
	}

	signpath := []string{}

	//遍历目录
	filepath.Walk(path,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			//匹配
			matched := reg.MatchString(f.Name())
			if matched {
				dirPath := filepath.Dir(path)
				signpath = append(signpath, dirPath)
			}
			return nil
		})
	return signpath
}

func (doc Doc) getDocPath(dirPath string) ([]string, error) {
	var files = []string{}

	dir, err := ioutil.ReadDir(dirPath)

	suffix := strings.ToUpper(doc.extension) //忽略后缀匹配的大小写

	if err != nil {
		return nil, err
	} else {
		for _, fi := range dir {
			if fi.IsDir() { // 忽略目录
				continue
			}
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
				files = append(files, dirPath+PthSep+fi.Name())
			}
		}
		return files, nil
	}

}

//遍历标记文件夹文件
func (doc Doc) outDoc(root string, outpath string) {
	//制定匹配模式
	//docMatch := `^` + doctab + `(.)+__end__$`
	//docMatch := `(?i:^//|/\*|/\*\*__doc__).*//|/\*|/\*\*__end__`
	//reg := regexp.MustCompile(docMatch)
	docpaths := getSignPath(root)
	for _, path := range docpaths {
		if filenames, err := doc.getDocPath(path); err == nil {
			for _, file := range filenames {
				ftext, _ := ioutil.ReadFile(file)
				//text := reg.FindAllString(string(ftext), -1)
				fmt.Printf("%s\n", ftext)
			}
		} else {
			fmt.Printf("遍历文档目录错误：%s", err)
		}

	}
}

func (doc Doc) General() {
	doc.outDoc(doc.root, doc.docpath)
	//s := getSignPath(doc.root)
	//fmt.Println(getSignPath(doc.root))
}
