package wp_json

import (
	"github.com/xiao-ha/request"
	"fmt"
	"encoding/json"
	"github.com/json-iterator/go"
)
/*
http://demo.wp-api.org/wp-json/wp/v2/users
*/

type wp_post struct {
	//publish, future, draft, pending, private
	Status	string	`json:"status"`
	Title	string	`json:"title"`
	Content	string	`json:"content"`
	Author	int		`json:"author"`
	Categories []int `json:"categories"`
	Tags	[]int	`json:"tags"`
}

type wp_tag struct{
	Description string		`json:"description"`
	Name		string		`json:"name"`
	Slug		string		`json:"slug"`
	Meta		interface{}	`json:"meta"`
}

type wp_category struct{
	Description string		`json:"description"`
	Name		string		`json:"name"`
	Slug		string		`json:"slug"`
	Parent		int			`json:"parent"`
	Meta		interface{}	`json:"meta"`
}

type list_wp_categories struct {
	Id			int			`json:"id"`
	Count		int			`json:"count"`
	Description	string		`json:"description"`
	Link		string		`json:"link"`
	Name		string		`json:"name"`
	Slug		string		`json:"slug"`
	Taxonomy	string		`json:"taxonomy"`
	Parent		int			`json:"parent"`
	Meta		interface{}	`json:"meta"`
	_links		interface{}	`json:"-"`
}

type list_wp_tags struct{
	Id			int			`json:"id"`
	Count		int			`json:"count"`
	Description string		`json:"description"`
	Link		string		`json:"link"`
	Name		string		`json:"name"`
	Slug		string		`json:"slug"`
	Taxonomy	string		`json:"taxonomy"`
	Meta		interface{}	`json:"meta"`
	_links		interface{}	`json:"-"`
}

const(
	PUBLISH = iota
	FUTURE
	DRAFT
	PENDING
	PRIVATE
)

type WP_JSON struct{
	web	*request.Request
	uri	string
	categories []list_wp_categories
	tags []list_wp_tags
}

///private
func (wj *WP_JSON)getCategoryIdBySlug(slug string) int{
	for _ , b:= range wj.categories{
		if b.Slug == slug{
			return b.Id
		}
	}
	return -1
}



func (wj *WP_JSON)getTagIdBySlug(slug string) int{
	for _ , b:= range wj.tags{
		if b.Slug == slug{
			return b.Id
		}
	}
	return -1
}

//public
func Build(uri string , username string ,password string) *WP_JSON {
	var wj WP_JSON
	wj.web = request.Build("" , username,password)
	wj.uri = uri
	return &wj
}

func (wj *WP_JSON)IsCategoryExist(slug string) bool{
	return wj.getCategoryIdBySlug(slug) != -1
}

func (wj *WP_JSON)IsTagExist(slug string) bool{
	return wj.getTagIdBySlug(slug) != -1
}

func (wj *WP_JSON)LoadCategories(){
	var r []list_wp_categories
	for i:=1 ; ; i++{
		uri := fmt.Sprintf("%s/wp-json/wp/v2/categories?page=%d", wj.uri  , i)
		k := wj.web.Get(uri , "")
		if len(k) > 2{
			var b []list_wp_categories
			err := json.Unmarshal([]byte(k), &b)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			fmt.Println(b)
			r = append(r , b...)
		}else{
			break
		}
	}
	wj.categories = r
}

func (wj *WP_JSON)LoadTags(){
	var r []list_wp_tags
	for i:=1 ; ; i++ {
		uri := fmt.Sprintf("%s/wp-json/wp/v2/tags?page=%d", wj.uri  , i)
		k := wj.web.Get(uri, "")
		if len(k) > 2 {
			var b []list_wp_tags
			err := json.Unmarshal([]byte(k), &b)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			fmt.Println(b)
			r = append(r , b...)
		}else{
			break
		}
	}
	wj.tags = r
}

func (wj *WP_JSON)CreateTag(name string , slug string , description string) (bSuccess bool,error_msg string) {
	data := wp_tag{
		Description:description,
		Name:name,
		Slug:slug,
	}
	uri := wj.uri + "/wp-json/wp/v2/tags"
	r := wj.web.Post(uri ,data , "" , request.JSON)
	if len(jsoniter.Get([]byte(r), "slug").ToString()) > 0{
		return true , ""
	}
	return false , r
}

func (wj *WP_JSON)CreateCategory(name string , slug string , description string , parent string) (bSuccess bool,error_msg string){
	data := wp_category{
		Description:description,
		Name:name,
		Slug:slug,
	}

	if len(parent) > 0 {
		nParent := wj.getCategoryIdBySlug(parent)
		if nParent != -1{
			data.Parent = nParent
		}
	}

	uri := wj.uri + "/wp-json/wp/v2/categories"
	r := wj.web.Post(uri ,data , "" , request.JSON)
	if len(jsoniter.Get([]byte(r), "slug").ToString()) > 0{
		return true , ""
	}
	return false , r
}


func (wj *WP_JSON)WritePost(title string ,content string , status int , sCategories []string , sTags []string) (bSuccess bool,error_msg string){
	nCategories := make([]int , len(sCategories))
	for i,c := range sCategories{
		nCategories[i] = wj.getCategoryIdBySlug(c)
		if nCategories[i] == -1{
			return false,"Some categories no exist!"
		}
	}

	nTags := make([]int , len(sTags))
	for i,c := range sTags{
		nTags[i] = wj.getTagIdBySlug(c)
		if nTags[i] == -1{
			return false,"Some tags no exist!"
		}
	}

	data := wp_post{
		Title:title,
		Content:content,
		Author:1,
		Categories:nCategories,
		Tags:nTags,
	}

	switch status {
	case PUBLISH:
		data.Status = "publish"
	case FUTURE:
		data.Status = "future"
	case DRAFT:
		data.Status = "draft"
	case PENDING:
		data.Status = "pending"
	case PRIVATE:
		data.Status = "private"
	default:
		data.Status = "draft"
	}

	uri := wj.uri + "/wp-json/wp/v2/posts"
	r := wj.web.Post(uri ,data , "" , request.JSON)
	if len(jsoniter.Get([]byte(r), "date").ToString()) > 0{
		return true , ""
	}
	return false , r
}