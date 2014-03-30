package main

import ( 
    "fmt"
    "math/rand"
    "time" 
    "encoding/json" 
    "errors"
    zmq "github.com/alecthomas/gozmq" 
)
type Request struct {
    Method string
    Args interface{}
} 

type Response struct {
    Response interface{}
    Status int
} 

func genrands(num int) (ret [] float64) { 

    ret = make([] float64, num)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < num; i++ {
        ret[i] = r.Float64() * 100
    }
    return ret
    
} 

func call(socket *zmq.Socket, method string, 
        args ...interface{}) (interface{}, error){ 
    
    msg := Request{method, args} 

    enc, err := json.Marshal(msg)
    if err != nil {
        fmt.Println("json error", err) 
        return -1, err        
    } 
    socket.Send(enc, 0)
    resp := new(Response) 
    
    // Wait for reply, could probably start a timeout thing
    reply, _ := socket.Recv(0)
        
    err = json.Unmarshal(reply,resp) 
    
    if err != nil {
        fmt.Println("json error", err)
        return -1, err
    }

    if resp.Status != 0 {
        return -1, errors.New("RPC Method not found:"+method) 
    } 
    
    return resp.Response, nil

}


func main() {

    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.REQ)
    defer context.Close()
    defer socket.Close()

    fmt.Println("Connecting to python server…")
    socket.Connect("tcp://localhost:5555")

    for i := 0; i < 10; i++ {
        floats := genrands(100)
        sum,_ := call(socket, "sum", floats)
        std,_ := call(socket, "std", floats)
        variance,_ := call(socket, "var", floats) 
        mean,_ := call(socket, "mean", floats) 
        
        add, _ := call(socket, "add", 2,3)

        fmt.Println(sum,std,variance,mean,add) 
    }


} 
