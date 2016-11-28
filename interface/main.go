// Goのインターフェースの理解をするうえで、下記のブログが丁寧で分かりやすかったため写経させてもらっています。
// See: http://jxck.hatenablog.com/entry/20130325/1364251563

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Run() {
	fmt.Println("<----------main1:---------->")
	main1()

	fmt.Println("<----------main2:---------->")
	main2()

	fmt.Println("<----------main3:---------->")
	main3()

	fmt.Println("<----------main4:---------->")
	main4()

	fmt.Println("<----------main5:---------->")
	main5()

	fmt.Println("<----------main6:---------->")
	main6()

	fmt.Println("<----------main7:---------->")
	main7()

	fmt.Println("<----------main8:---------->")
	main8()

	fmt.Println("<----------main9:---------->")
	main9()

	fmt.Println("<----------main10:---------->")
	main10()

	fmt.Println("<----------main11:---------->")
	main11()

	fmt.Println("<----------main12:---------->")
	main12()

	fmt.Println("<----------main13:---------->")
	main13()

	fmt.Println("<----------main14:---------->")
	main14()
}

// 基本的な struct
type Point struct {
	X int
	Y int
}

func main1() {
	var a Point = Point{2, 3}
	fmt.Println(a.Coordinate())
}

func (p Point) Coordinate() string {
	// p がレシーバー
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// 既存の型を元にした独自の型を定義できる。
type trimmedString string

func (t trimmedString) trim() trimmedString {
	return t[:3]
}

func main2() {
	var t trimmedString = "abcdefg"
	fmt.Println(t.trim())

	var s string = string(t)
	fmt.Println(s)
}

//// Interface を宣言
type Accessor interface {
	GetText() string
	SetText(text string)
}

//// Accessor を満たす実装
//// Interface の持つメソッド群を実装していれば、
//// Interface を満たす(satisfy) といえる。
//// 明示的な宣言は必要なく、実装と完全に分離している。
type Document struct {
	text string
}

func (d *Document) GetText() string {
	return d.text
}

func (d *Document) SetText(text string) {
	d.text = text
}

func main3() {
	doc := &Document{}
	doc.SetText("document")
	fmt.Println(doc.GetText())

	// Accessor Interface を実装しているので
	// Accessor 型に代入可能
	var accessor Accessor = &Document{}
	accessor.SetText("Accessor")
	fmt.Println(accessor.GetText())
}

type Page struct {
	Document //　匿名性を含むと、その型のメソッドが継承（というか mixin）される
	Page     int
}

func main4() {
	// Page は Document を継承しており
	// Accessor Interface を満たす。
	// この場合代入可能
	var accessor Accessor = &Page{}
	// この値は accessor.Document.text に設定されている。
	// accessor の構造体がレシーバーになっているわけではないということ
	accessor.SetText("page")
	fmt.Println(accessor.GetText())

	// Document と Page　の間に代入可能な関係は無い
	// var page Page = Document{}
	// var doc Document = Page{}
}

/*
	Duck Typing
	Accessor を満たしていれば、 Get, Set できるという例。
*/
func SetAndGet(accessor Accessor) {
	accessor.SetText("accessor")
	fmt.Println(accessor.GetText())
}

func main5() {
	// どちらも Accessor として振る舞える
	SetAndGet(&Page{})
	SetAndGet(&Document{})
}

/*
	Override
*/
type ExtendedPage struct {
	Document
	Page int
}

// Document.GetText() のオーバーライド
func (ep *ExtendedPage) GetText() string {
	// int -> string は strconv.Itoa 使用
	return strconv.Itoa(ep.Page) + " : " + ep.Document.GetText()
}

func main6() {
	// Accessor を実装している
	var accessor Accessor = &ExtendedPage{
		Document{},
		2,
	}
	accessor.SetText("page")
	fmt.Println(accessor.GetText()) // 2 : page
}

// Interface 型
// er を付ける命名が習慣
type Getter interface {
	GetText() string
}

func dynamicIf(v interface{}) string {
	// v はInterface型
	var result string
	g, ok := v.(Getter) // v が Get() を実装しているか調べる
	if ok {
		result = g.GetText()
	} else {
		result = "not implemented"
	}
	return result
}

func dynamicSwitch(v interface{}) string {
	// v は Interface 型

	var result string

	// v が実装している型でスイッチする
	switch checked := v.(type) {
	case Getter:
		result = checked.GetText()
	case string:
		result = "not implemented"
	}
	return result
}

func main7() {
	var ep *ExtendedPage = &ExtendedPage{
		Document{},
		3,
	}
	ep.SetText("page")

	// do は Interface　型を取り
	// ジェネリクス的なことができる
	// 型スイッチを行う場合
	fmt.Println(dynamicIf(ep))           // 3: page
	fmt.Println(dynamicSwitch("string")) // not implemented
}

// 全ての型を許容するインターフェースのようなものを作っておく
type Any interface{}

// ジェネリクス的な
type GetValuer interface {
	GetValue() Any
}

// Any型で実装
type Value struct {
	v Any
}

// GetValuer を実装
func (v *Value) GetValue() Any {
	return v.v
}

func main8() {
	// インターフェースで受け取る
	var i GetValuer = &Value{10}
	var s GetValuer = &Value{"vvv"}

	// インターフェース型のコレクションに格納
	var values []GetValuer = []GetValuer{i, s}

	// それぞれ GetValue()　が Any で呼べる
	for _, val := range values {
		fmt.Println(val.GetValue())
	}
}

func PrintAll(values []interface{}) {
	for _, val := range values {
		fmt.Println(val)
	}
}

func main9() {
	names := []string{"one", "two", "three"}

	// これは間違い
	// PrintAll(names)

	// 明示的に変換が必要
	values := make([]interface{}, len(names))
	for i, v := range names {
		values[i] = v
	}
	PrintAll(values)
}

// 1, twitter API から Time のパース

/*
	Twitter の JSON を map にパースする。
	twitter の JSON には、時間が Ruby フォーマットの文字列で格納されているので、
	それを考慮して型を考える。
	"Thu May 31 00:00:01 +0000 2012"
*/

var JSONString = `{ "created_at": "Thu May 31 00:00:01 +0000 2012" } `

func main10() {
	// map として {string: interface{}} としてしまえば
	// value　がなんであれパースは可能
	var parsedMap map[string]interface{}

	if err := json.Unmarshal([]byte(JSONString), &parsedMap); err != nil {
		panic(err)
	}

	fmt.Println(parsedMap) // map[created_at:Thu May 31 00:00:01 +0000 2012]
	for k, v := range parsedMap {
		fmt.Println(k, reflect.TypeOf(v)) // created_at string
	}
}

type Timestamp time.Time

// Unmarshaller　を実装
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	fmt.Println("UnmarshalJSON")
	v, err := time.Parse(time.RubyDate, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = Timestamp(v)
	return nil
}

func main11() {
	var val map[string]Timestamp // 定義した型を使う

	if err := json.Unmarshal([]byte(JSONString), &val); err != nil {
		panic(err)
	}

	// パースされていることを確認
	for k, v := range val {
		fmt.Println(k, time.Time(v), reflect.TypeOf(v))
		// created_at 2012-05-31 00:00:01 +0000 +0000 main.Timestamp
	}
}

// 各型が、自身のパース実装を持てばよいので、そのメソッドだけ定義しておく。
type Entity interface {
	UnmarshallJSON([]byte) error
}

func GetEntity(b []byte, e Entity) error {
	// 各実装に処理を移譲
	return e.UnmarshallJSON(b)
}

// 型を定義
// User に関する必要なデータだけ取りたい型的な
type UserData struct {
	Id        int
	Name      string
	Time_Zone string
	Lang      string
}

// *_count だけ適当に取りたい型的な
type CountData struct {
	Followers_count  int
	Friends_count    int
	Listed_count     int
	Favourites_count int
	Statuses_count   int
}

// Entity を実装
// ここでは、 json モジュールになげるだけで
// 同じ実装でできてしまったが、
// 本来 Entity ごとに違う実装になる。
func (d *UserData) UnmarshallJSON(b []byte) error {
	err := json.Unmarshal(b, d)
	if err != nil {
		return err
	}
	return nil
}

func (d *CountData) UnmarshallJSON(b []byte) error {
	err := json.Unmarshal(b, d)
	if err != nil {
		return err
	}
	return nil
}

func main12() {
	// 対象の JSON 文字列
	EntityString := `{
		"id":51442629,
		"name":"Jxck",
		"followers_count":1620,
		"friends_count":617,
		"listed_count":204,
		"favourites_count":2895,
		"time_zone":"Tokyo",
		"statuses_count":17387,
		"lang":"ja"
	}`
	userData := &UserData{}
	countData := &CountData{}
	GetEntity([]byte(EntityString), userData)
	GetEntity([]byte(EntityString), countData)
	fmt.Println(*userData)  // {51442629 Jxck Tokyo ja}
	fmt.Println(*countData) // {1620 617 204 2895 17387}
}

// タグ付きの struct を定義
type TaggedStruct struct {
	field string `tag:"tag1"`
}

func main13() {
	// reflect　でタグを取得
	var ts = TaggedStruct{}
	var t reflect.Type = reflect.TypeOf(ts)
	var f reflect.StructField = t.Field(0)
	var tag reflect.StructTag = f.Tag
	var val string = tag.Get("json")
	fmt.Println(tag, val)
}

// JSON をマッピングするために
// キー名のタグをつけた struct を定義
type Employee struct {
	Name  string `json:"emp_name"`
	Email string `json:"emp_email"`
	Dept  string `json:"dept"`
}

func main14() {
	// フィールド名が struct の Filed　名と違う JSON も
	// json:"fieldName" の形でタグを付けてあるので
	// マッピングすることができる。
	var jsonString = []byte(`
	{
		"emp_name" : "john",
		"emp_email": "john@golang.com",
		"dept" : "HR"
	}`)

	var john Employee
	err := json.Unmarshal(jsonString, &john)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", john) // {Name:john Email:john@golang.com Dept:HR}
}
