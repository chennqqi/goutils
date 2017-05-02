#closeevent



##demo


func demo(){
	c := make(chan os.Signal)
	signalNotify(c)
	for {
		fmt.Println("run...", err)
		
		s := <-c

		//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
		fmt.Println("get signal:", s)
		fmt.Println("close...")
		
		//do your own close...
		//doCloseFunction()
			
		break
	}
}