# autodomain
动态IP自动同步到万网域名

使用说明
附systemctl脚本<自行修改>:
/*************************************/
[Unit]  
Description=autodomain  
After=network.target  
      
[Service]  
Type=simple  
User=yiye
ExecStart=/opt/autodomain
PrivateTmp=true
      
[Install]  
WantedBy=multi-user.target
/*************************************/
/*************************************/
配置config.json
{
    "Time" : 120,  //循环时间单位为秒
    "Id" : "",    //万网提供的access_key_id
    "Secret" :"",  //万网提供的access_key_secret
    "Domain":""   //要修改的域名 支持的类型为: A
    "CheckUrls": ["http://ddns.oray.com/checkip","http://www.whatismyip.com.tw/"] //检测IP的网站，干扰信息越少越好。
} 
/*************************************/
/*************************************/
自己的VPS可使用以下php代码自检IP检测页面。以下代码保存为.php文件放入www目录下。
php code :
/********************php代码开始，不包含此行************************/
<?php
function getIP()
{
    static $realip;
    if (isset($_SERVER)){
        if (isset($_SERVER["HTTP_X_FORWARDED_FOR"])){
            $realip = $_SERVER["HTTP_X_FORWARDED_FOR"];
        } else if (isset($_SERVER["HTTP_CLIENT_IP"])) {
            $realip = $_SERVER["HTTP_CLIENT_IP"];
        } else {
            $realip = $_SERVER["REMOTE_ADDR"];
        }
    } else {
        if (getenv("HTTP_X_FORWARDED_FOR")){
            $realip = getenv("HTTP_X_FORWARDED_FOR");
        } else if (getenv("HTTP_CLIENT_IP")) {
            $realip = getenv("HTTP_CLIENT_IP");
        } else {
            $realip = getenv("REMOTE_ADDR");
        }
    }
    return $realip;
}
echo getIP();
?>
/********************php代码结束，不包含此行************************/

v4.
1.Time属性已可用默认为秒，内部循环方式运行。
2.加入日志显示，但并未加入日志输出。
3.config.json文件放入HOME目录下
4.新增多checkip地址检测
5.修复部分BUG

v5.
1.修正多checkip地址中不返回IP的BUG.
2.修正部分提示信息.
