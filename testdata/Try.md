> 从源码安装  
你可以通过以下两种方式下载安装 goconfig：  
从Gopm快速安装  
```
gopm get github.com/Unknwon/goconfig
```
或从Github快速安装  
```
go get github.com/Unknwon/goconfig
```

> 导入包
```
"github.com/Unknwon/goconfig"
```

> 创建ConfigFile对象 对ini文件的操作会用到 
```
//支持绝对路径
filename := ""
//返回一个对象和一个错误
cfg, err := goconfig.LoadConfigFile(filename)
if err != nil{//这里错误处理}
```
> 加入更多的ini文件
```
// More File Test
	path := "src/go-config/"
	fileName1 := path + "config.ini"
	fileName2 := path + "conf1.ini"
	fileName3 := path + "conf2.ini"
	cfg, err := goconfig.LoadConfigFile(fileName1, fileName2) // 在这里可以添加更多的文件
	if err != nil {
		log.Fatalf("INI File %s", err)
	}
	log.Print("INI File Open Success!")
```

> 取得指定分区指定键的值
func GetValue
e.g :
```
	section := "" //DEFAULT分区
	key := "key1" //分区下的Key值
	req, err := cfg.GetValue(section, key) //此处的cfg为你的ConfigFile对象
    if err != nil {//这里错误处理}
    //简写(前提是分区和key必须找到到 错误返回空)：
    req = cfg.MustValue(section,key)
```

> 设置指定分区指定键的值
e.g :
```
    key := "key1" //分区下的Key值
    value := "New Value"
	if cfg.SetValue("",//DEFAULT分区
		key, value) {
		log.Printf("Set/Delect %s Success! %s ", key, value)
		return
	}
```

> 取得指定分区指定键的注释 func GetSectionComments
```
    section := "" //DEFAULT分区
	comment := cfg.GetSectionComments(section) //如果出现无法找到的错误 返回空
	log.Printf("分区 %s 的注释 %s", section, comment) //Print
```

> 设置指定分区指定键的注释 func SetSectionComments
// 如果该位置已经存在注释 直接覆盖返回false 如果不存在 直接插入返回ture
```
    contents := "# 分区注释会在这里"
	if cfg.SetSectionComments(section,
		contents){
			log.Printf("Set/Delect %s Comment Success! %s ",
				section,contents)
	}
	log.Printf("Update %s Comment Success! %s ",
		section, contents)
```  
> 取得指定分区下的全部键值对 如果不存在 返回nil和错误
```
// Get All Section Contents
	section = "auto increment" //分区
	sec,err := cfg.GetSection(section)
	if err != nil {
		log.Printf(
			"No Such %s -> %s",
			section, err)
	}
	log.Printf("%s -> %s ",
		section, sec)
```
--- 如果你想手动指定的返回类型 以下的方法会很适合你 你甚至可以指定缺省值{默认值}
func Int //返回int类型
func Int64 //返回int64类型
func MustBool //返回bool类型  
e.g :
```
// Request Type Int
	section = "parent"
	key = "age"
	v1,err := cfg.Int(section,key)
	if err != nil {
		log.Printf("Section %s EOF ! %s",section,err)
	}
	log.Printf("%s -> %s: %d",
		section, key, v1)
	// Request Type Bool
	section = "parent.child"
	key = "married"
	v2 := cfg.MustBool(section,key)
	log.Printf("%s -> %s: %v",
		section, key, v2)
```
---想动态加载你的ini文件吗 支持的
```
// INI File Append... ...
	err = cfg.AppendFiles(fileName3) //这里添加你的ini文件
	if err!=nil {
		log.Fatalf("INI File Append %s",err)
	}
	log.Print("INI File Append Success!")
	//INI File Reload... ... 
	err = cfg.Reload() //重新加载配置文件以防其发生变化
	if err != nil {
		log.Print("INI File Reload EOF!")
	}
	log.Print("INI File Reload Success!")
```
> 子孙关系
```
// Request Type Int | Set Default Value | Father Son Test
	section = "parent.child" //分区
	key = "age" //Key
	v2 := cfg.MustInt(section,key,0) //指定返回的类型为Int 并 设置默认返回值 如果没有找到对应的Key就会返回
	log.Printf("%s -> %s: %v",
		section, key, v2)
```  
--- 你想删除你的key吗 请试试下面的方法  
func DeleteKey //如果该键已被删除，则返回true;如果该键或键不存在，则返回false  
--- 最后, 让我们来保存吧  
> 保存
```
// Saver
	filename = "src/go-config/conf1.ini"
	err = goconfig.SaveConfigFile(
		cfg,filename)
	if err != nil {
		log.Printf("Saver %s %s",filename,err)
	}
	log.Printf("Saver %s Success!",filename)
}
```





