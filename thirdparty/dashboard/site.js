var app = new Vue({
  el: '#app',
  data: {
    isLoading:false,
    taskList: [
      { taskName: "taskOne", taskType:"OneShot", taskState:"Running",lastExecutedTime:"2017/05/16 15:00:46",lastExecuteCost:"50s"},
      { taskName: "taskTwo", taskType:"Plan", taskState:"Executing",lastExecutedTime:"2017/05/16 15:00:46",lastExecuteCost:"50s"},
      { taskName: "taskThree", taskType:"Plan", taskState:"Stopped",lastExecutedTime:"2017/05/16 15:00:46",lastExecuteCost:"50s"},
      { taskName: "taskFour", taskType:"OneShot", taskState:"Completed",lastExecutedTime:"2017/05/16 15:00:46",lastExecuteCost:"50s"},
    ],
    curTask: {},
  },
  created() {
    this.refreshList()

    setInterval(()=>{
      this.refreshList()
    }, 5000)
  },
  methods:{
    refreshList(){
      this.simpleAjax("/GetAllTasks").then((allTasks)=>{
        this.taskList = allTasks
        if (!this.curTask.taskName){
          this.curTask = this.taskList[0]
        }else{
          this.curTask = this.taskList.filter(p=>p.taskName == this.curTask.taskName)[0]
        }
      })
    },
    clickTask(aTask){
      this.curTask = aTask
    },
    clickStart(ev){
      if(this.isLoading){
        return;
      }
      this.isLoading = true;
      this.simpleAjax("/StartTask?taskname="+ this.curTask.taskName).then(()=>{
        swal({
          title: "Success!",
          text: "Every is fine, you have started a task just now!",
          icon: "success",
        });
        this.refreshList();
      }).catch((rs)=>{
        swal({
          title: "Error!",
          text: "Something wrong,code:" + rs.code + ",msg:"+ rs.msg,
          icon: "error",
        })
      })
    },
    clickExecute(ev){
      if(this.isLoading){
        return;
      }
      this.isLoading = true;
      this.simpleAjax("/ExecuteTask?taskname="+ this.curTask.taskName).then(()=>{
        swal({
          title: "Success!",
          text: "Every is fine, you have execueted a task just now!",
          icon: "success",
        });
        this.refreshList();
      }).catch((rs)=>{
        swal({
          title: "Error!",
          text: "Something wrong,code:" + rs.code + ",msg:"+ rs.msg,
          icon: "error",
        })
      })
    },
    clickStop(ev){
      if(this.isLoading){
        return;
      }
      this.isLoading = true;
      this.simpleAjax("/StopTask?taskname="+ this.curTask.taskName).then(()=>{
        swal({
          title: "Success!",
          text: "Every is fine, you have stopped a task just now!",
          icon: "success",
        });
        this.refreshList();
      }).catch((rs)=>{
        swal({
          title: "Error!",
          text: "Something wrong,code:" + rs.code + ",msg:"+ rs.msg,
          icon: "error",
        })
      })
    },
    clickReset(ev){
      if(this.isLoading){
        return;
      }
      this.isLoading = true;
      this.simpleAjax("/ResetTask?taskname="+ this.curTask.taskName).then(()=>{
        swal({
          title: "Success!",
          text: "Every is fine, you have reseted a task just now!",
          icon: "success",
        });
        this.refreshList();
      }).catch((rs)=>{
        swal({
          title: "Error!",
          text: "Something wrong,code:" + rs.code + ",msg:"+ rs.msg,
          icon: "error",
        })
      })
    },
    simpleAjax(url){
      return new Promise((rs,rj)=>{
        var hr =new XMLHttpRequest()
        hr.open('get', url, true)
        hr.responseType = 'json'
        hr.onreadystatechange = ()=>{
          if(hr.readyState == 4){
            if(hr.status == 200){
              rs(hr.response)
            } else {
              rj({ code:hr.status, msg:hr.response.msg })
            }
          }
          this.isLoading = false;
        }
        hr.send()
      })
    }
  }
})