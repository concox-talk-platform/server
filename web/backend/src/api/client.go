/**
* @Author: yanKoo
* @Date: 2019/3/20 16:15
* @Description:
 */
package api

//const Port = "9000"

//func main() {
//	conn, err := grpc.Dial(":"+Port, grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("grpc.Dial err: %v ", err)
//	}
//
//	defer func() {
//		if err := conn.Close(); err != nil {
//			log.Fatalf("connection err: %v", err)
//		}
//	}()
//
//	client := pb.NewAccountServiceClient(conn)
//	resp, err := client.SignOut(context.Background(), &pb.SignOutRequest{
//		SessionId: "1234",
//	})
//	if err != nil {
//		log.Fatalf("client SignOut err : %v", err)
//	}
//
//	log.Println(resp.GetResult())
//}
