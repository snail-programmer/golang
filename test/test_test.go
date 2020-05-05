package testfuck
import(
	"testing"
)
// func Test_Div_1(t *testing.T){
// 	if r, e := Div(6, 2); r!=3 || e !=nil{
// 		t.Error("测试未通过")	
// 	}else{
// 		t.Log("you are passed!")
// 	}
// }
// func Test_Div_2(t *testing.T){
// 	t.Error("测试真的未通过")
// }
func Benchmark_Division(b *testing.B) {
    for i := 0; i < b.N; i++ { //use b.N for looping 
        Div(4, 5)
    }
}
func Benchmark_TimeConsumingFunction(b *testing.B) {
    b.StopTimer() 
      b.StartTimer() //重新开始时间
    for i := 0; i < b.N; i++ {
        Div(4, 5)
    }
}