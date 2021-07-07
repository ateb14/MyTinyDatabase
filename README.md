# MyTinyDatabase. 
This is my first try in a Golang project.(Maybe in my coding journey, too). 
Just for fun lol. 
数据的储存单元为Unit，以规定的格式逐条写入tinydb.data中以方便编码和解码,每个Unit通过hash映射得到Offset，代表这个Unit在文件中出现的位置。  
文件的增删通过日志更新进行,Unit中具有Tag表示状态，PUT代表已加入文件，DEL代表已删除文件，为了防止冗余数据过多，通过Update操作清楚冗余数据（目前尚未测试，但绝对有bug >_<)  
目前还有不少bug没有修，不过接下来两天有些事，所以就先丢在这了...  
