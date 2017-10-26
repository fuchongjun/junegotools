package study

import (
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/axgle/mahonia"
	"runtime"
	"time"
	"encoding/json"
)
type PacsInfo struct {
	PID      string `json:"pid"`
	NC       string `json:"name"`
	SX       string `json:"sex"`
	BR       time.Time `json:"birthday"`
	Modality string `json:"modality"`
	DISKID   string `json:"picpath"`
}
func OdbcTestAccess() {
	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=E:\\June\\WorkSpace\\Pacs\\PACS.MDB")
	//conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb)};dbq=E:\\June\\WorkSpace\\Pacs\\PACS.MDB")//32位系统

	fmt.Println(runtime.GOARCH, runtime.GOOS)

	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	defer conn.Close()
	stmt, err := conn.Prepare("select A.PID,A.NC,A.SX,A.BR,B.Modality,B.DISKID from PATIENT A INNER JOIN STUDY B ON A.PID=B.PID WHERE A.PID='00022236'")
	if err != nil {
		fmt.Println("Query Error", err)
		return
	}
	defer stmt.Close()
	row, err := stmt.Query()
	if err != nil {
		fmt.Println("Query Error", err)
		return
	}
	defer row.Close()
	for row.Next() {
		pacsInfo :=new(PacsInfo)
		if err := row.Scan(&pacsInfo.PID, &pacsInfo.NC,&pacsInfo.SX, &pacsInfo.BR, &pacsInfo.Modality, &pacsInfo.DISKID); err != nil {
			fmt.Println(err)
		}
		decoder := mahonia.NewDecoder("gb18030")
		pacsInfo.NC=decoder.ConvertString(pacsInfo.NC)//gbk转为utf8
		pacsInfo.SX=decoder.ConvertString(pacsInfo.SX)//gbk转为utf8

		pacsinJsons,err:=json.Marshal(pacsInfo)
		if err!=nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(pacsinJsons))


	}
	fmt.Printf("%s\n", "finish")
	return
}
