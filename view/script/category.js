     //分类对象-学科分类
    class Category{
        //name
        //contain
    }
    var cate = new Array();
    var gory=new Category();
    var goryIx=0;
    var cateIx=0;
    gory.contain=new Array();
    //它把一个数组分解成右边的结构形式[obj1{name,contain[]}, obj2, ...]
     function category_callback(arr){
        gory.name=arr[0].CategoryName;
        $.each(arr,function(index,e){
            if(gory.name == e.CategoryName){
                gory.contain[goryIx++]=e.CategoryContain;
            }else{
                goryIx=0;
                cate[cateIx++]=gory;
                gory=new Category();
                gory.name=e.CategoryName;
                gory.contain=new Array();
                gory.contain[goryIx++]=e.CategoryContain;  
            }
        });
        cate[cateIx]=gory;
        goryIx=cateIx=0;
       console.log("cate-src:",cate,cate.length);
       if(this.cate_call){
           this.cate_call(cate);
       }
    }