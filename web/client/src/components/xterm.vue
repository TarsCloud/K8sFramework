<template>
  <div id="xterm" class="xterm-container"></div>
</template>

<script>
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";
import { WebLinksAddon } from "xterm-addon-web-links";

export default {
  data() {
    return {
      term: "",
      socket: "",
      app: "",
      server: "",
      pod: "",
      history: "",
      nodeip: "",
    };
  },
  mounted() {
    let { app, server, pod, history, nodeip } = this;
    app = this.getQueryVariable("AppName");
    server = this.getQueryVariable("ServerName");
    pod = this.getQueryVariable("PodName");
    history = this.getQueryVariable("History");
    nodeip = this.getQueryVariable("NodeIP");

    if (app === false || server === false || pod === false) {
      return alert("无法获取到容器，请联系管理员");
    }

    this.app = app;
    this.server = server;
    this.pod = pod;

    // 直接发给node后端，由node后端代理
    this.$ajax.getJSON("/server/api/shell_domain").then(domain => {
      let path = this.$ajax.ServerUrl.get(this.$ajax.parseUrl('/server/api/shell', { History: history, NodeIP: nodeip, AppName: app, ServerName: server, PodName: pod }))
      let url = `${domain.fromHost}${path}`
      console.info(url)
      this.init(url);
    })

  },
  methods: {
    //节流,避免拖动时候频繁向后端请求更新
    debounce(fn, wait) {
      let timeout = null;
      return function () {
        if (timeout !== null) clearTimeout(timeout);
        timeout = setTimeout(fn, wait);
      };
    },
    //页面重新resize的时候,需要重新告诉后端cols和rows
    resizeScreen() {
      const fitAddon = new FitAddon();
      this.term.loadAddon(fitAddon);
      fitAddon.fit();
      this.send(
        JSON.stringify({
          operation: "resize",
          cols: Math.floor(this.term.cols),
          rows: Math.floor(this.term.rows),
        })
      );
    },
    getQueryVariable(variable) {
      let query = window.location.search.substring(1);
      let vars = query.split("&");
      for (let i = 0; i < vars.length; i++) {
        let pair = vars[i].split("=");
        if (pair[0] == variable) {
          return pair[1];
        }
      }
      return false;
    },
    initXterm() {
      this.term = new Terminal({
        rendererType: "canvas", //渲染类型
        convertEol: true, //启用时，光标将设置为下一行的开头
        scrollback: 30, //终端中的回滚量
        disableStdin: false, //是否应禁用输入
        cursorStyle: "block", //光标样式
        cursorBlink: true, //光标闪烁
      });

      //在绑定的组件上初始化窗口
      this.term.open(document.getElementById("xterm"));

      //全屏
      this.resizeScreen()

      //加载weblink组件
      this.term.loadAddon(new WebLinksAddon());

      //监听resize,当窗口拖动的时候,监听事件,实时改变xterm窗口
      window.addEventListener(
        "resize",
        this.debounce(this.resizeScreen, 1000),
        false
      );
      
      //聚焦
      this.term.focus();

      // 支持输入与粘贴方法
      let _this = this; //一定要重新定义一个this，不然this指向会出问题
      this.term.onData(function (key) {
        //这里key值是你输入的值，数据格式一定要找后端要！！！！
        let order = { operation: "stdin", data: key };
        _this.socket.onsend(JSON.stringify(order)); //转换为字符串
      });
    },
    init(url) {
      // 实例化socket
      this.socket = new WebSocket(url);
      // 监听socket连接
      this.socket.onopen = this.open;
      // 监听socket错误信息
      this.socket.onerror = this.error;
      // 监听socket消息
      this.socket.onmessage = this.getMessage;
      // 发送socket消息
      this.socket.onsend = this.send;
    },
    open: function () {
      this.initXterm();
      this.term.writeln(`connecting to pod \x1B[1;3;31m ${this.pod} \x1B[0m ... \r\n`);
    },
    error: function () {
      console.log("[error] Connection error");
    },
    close: function () {
      this.socket.close();
      console.log("[close] Connection closed cleanly");
      term.writeln("");
      window.removeEventListener("resize", this.resizeScreen);
    },
    getMessage: function (msg) {
      const data = msg.data && JSON.parse(msg.data);
      if (data.operation === "stdout") {
        this.term.write(data.data)
      }
    },
    send: function (order) {
      this.socket.send(order);
    },
  },
};
</script>

<style>
body {
  margin: 0;
  overflow: hidden;
  padding: 0;
}

.xterm-container {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
}
</style>
