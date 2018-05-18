package main

import (
	"github.com/xiao-ha/wp-json"
	"fmt"
)

func main(){
	r := wp_json.Build("http://blog.test.com" , "admin","LzdH hQ64 28wp XPRY NIuw GvIj")
	r.LoadCategories()
	r.LoadTags()
	r.LoadUsers()
	fmt.Println(r.IsUserExist("admin"))
	//k , m := r.CreateTag("tag5" , "tag5","tag5")
	//fmt.Println(k,m)
	//k , m := r.CreateCategory("category1" , "category1","category1" ,"qa")
	//fmt.Println(k,m)
	//k , m := r.WritePost("admin","test" , "test" , wp_json.DRAFT , []string{"qa"}, []string{"tag1","tag2"})
	//fmt.Println(k,m)
	fmt.Println(r.IsCategoryExist("qa"))
	fmt.Println(r.IsCategoryExist("666"))
	fmt.Println(r.IsTagExist("tag1"))
	fmt.Println(r.IsTagExist("666"))
}
