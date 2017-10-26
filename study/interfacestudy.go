package study

import "fmt"

type Connecter interface {
	Connect()
}
type USB interface {
	Name() string
	Connecter
}
type ComputerConnecter struct {
	ConnName string
}

func (computerconn ComputerConnecter) Name() string {
	return computerconn.ConnName
}
func (compuerconn ComputerConnecter) Connect() {
	fmt.Printf("链接成功!,%s\n", compuerconn.ConnName)
}
func getConnecterName(usb USB) string {
	if v, ok := usb.(ComputerConnecter); ok {
		return v.ConnName
	}
	return "未知接口类型"
}

func disconnect(usb interface{}) string {
	if v, ok := usb.(ComputerConnecter); ok {
		return "断开" + v.ConnName
	}
	return "未知接口类型"
}
func disconnectswich(usb interface{}) string {
	switch v := usb.(type) {
	case ComputerConnecter:
		return "断开" + v.ConnName
	default:
		return "未知接口类型"
	}
}

//接口类型转换，向上转换，子类转换父类可以
func ConvertInterface(comnn ComputerConnecter) Connecter {
	return Connecter(comnn)
}
func Test1() {
	var computer ComputerConnecter
	computer.ConnName = "电脑usb"
	fmt.Println(computer.Name())
	computer.Connect()

	fmt.Println(getConnecterName(computer))

	var conn = ConvertInterface(computer) //接口之间的赋值只是值拷贝
	conn.Connect()
	computer.ConnName = "接口名字更改"
	conn.Connect()
}

//=========学习flag.value 解析参数======
type Celsius float64
func (c Celsius)String()string  {
	return fmt.Sprintf("%g℃",c)
}


type celsiusFlag struct {
	Celsius
}

func (f *celsiusFlag)Set(s string)error {

	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}