# 1、fping

```
┌──(kali㉿kali)-[~]
└─$ fping -h                                                                             
Usage: fping [options] [targets...]

探测选项:
   -4, --ipv4         只 ping IPv4 地址
   -6, --ipv6         只 ping IPv6 地址
   -b, --size=BYTES   自定义要发送的ping的数据包大小，以字节为单位（默认值：56）
   -B, --backoff=N    设置指数补偿系数为N（默认值：1.5，范围1.0—5.0之间）
   -c, --count=N      计数模式：向每个目标发送N次ping
   -d, --rdns         使用DNS查找返回ping数据包的地址。这使您可以为fping提供IP地址列表作为输入，并在                       输出中显示主机名。这类似于选项-n / -name，但是即使您将主机名作为目标（NAME->                         IP-> NAME），也会强制执行反向DNS查找。
   -f, --file=FILE    从文件中读取目标列表（-表示标准输入）
   -g, --generate     生成目标IP列表 (仅当未指定-f时使用)
                      (给出目标列表的开始和结束IP地址，或者CIDR地址)
                      (例. fping -g 192.168.1.0 192.168.1.255 或 fping -g 192.168.1.0/24)
   -H, --ttl=N        设置IP的TTL值（Time To Live hops：生存时间跳数）
   -I, --iface=IFACE  指定特定网卡ping
   -l, --loop         循环模式：一直发送ping
   -m, --all          向目标主机的每一个IP地址发送ping（包括IPv4和IPv6），与-A一起使用
   -M, --dontfrag     设置IP标头中的“不分片”位（用于确定/测试MTU）
   -O, --tos=N        在ICMP数据包上设置服务类型（tos），N可以是十进制或十六进制（0xh）格式
   -p, --period=MSEC  设置ping数据包到一个目标的时间间隔（单位：毫秒）
                      (在循环和计数模式下，默认值：1000ms)
   -r, --retry=N      ping重试次数 (默认值: 3)
   -R, --random       随机分组数据(为了阻止链路数据压缩)，代替全0作为分组数据，将ping生成随机字节，                         来阻止像链路数据压缩的情形。
   -S, --src=IP       设置源IP地址
   -t, --timeout=MSEC 设置ping到单个目标IP初始超时时间。 (默认: 500 ms,
                      但 -l/-c/-C除外, 其中-p周期最长为2000ms)

输出选项:
   -a, --alive        显示存活的主机
   -A, --addr         显示目标地址
   -C, --vcount=N     与-c相同，报告以详细格式结果
   -D, --timestamp    在每个输出行之前打印时间戳
   -e, --elapsed      显示返回数据包经过的时间
   -i, --interval=MSEC  自定义发送ping报文的时间间隔(默认为10ms)
   -n, --name         显示目标主机名（与-d等效）
   -N, --netdata      与netdata兼容的输出（需要-l -Q）
   -o, --outage       显示累计中断时间（丢失的数据包/报文时间间隔）
   -q, --quiet        安静模式（不显示按目标或者按ping的结果）
   -Q, --squiet=SECS  与-q相同，但是每n秒显示一次摘要
   -s, --stats        打印最终统计
   -u, --unreach      显示无法达到的目标
   -v, --version      显示fping版本
   -x, --reachable=N  显示> = N个主机是否可访问
```
