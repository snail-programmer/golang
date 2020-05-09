package DBModel

type User struct {
	Id          string
	PhoneNumber string
	NickName    string
	Password    string
	School      string
	Education   string
	Description string
	HeadImg     string
	Coin        string
	Gratuity    string
	Money       string
	//virtual field
	Note_sum string
	View_sum string
	Col_sum  string
	Flow     string //流量收益
}
type Draft struct {
	Id         string
	AuthorId   string
	DraftNote  string
	CreateTime string
}
type Article struct {
	ArticleId  string
	AuthorId   string
	Article    string
	View_num   string
	Collection string
	Remark     string
	CreateTime string
	CategoryId string
	//virtual field
	NickName        string
	CategoryName    string
	CategoryContain string
	HeadImg         string
}
type Article_comment struct {
	Id        string
	ArticleId string
	Comment   string
}
type Category struct {
	Id              string
	CategoryName    string
	CategoryContain string
	CateType        string
}
type Black_list struct {
	Id       int
	MyId     int
	AuthorId int
}
type Collect_article struct {
	Id        string
	MyId      string
	ArticleId string
}
type Visit_log struct {
	Id         string
	VisitId    string
	CategoryId string
}
type Reg_log struct {
	Id string
}
