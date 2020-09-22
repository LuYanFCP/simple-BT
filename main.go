package main
import (
	"fmt"
	"simple-BT/torrent-client"
)

func main() {
	tf, err := BTClient.OpenTorrentFile("test/[BYRBT].VMware-workstation-full-16.0.0-16894299.exe.torrent")
	if err != nil {
		fmt.Println("无法打开，没有文件或者路径错误!")
	}
	err = tf.Download("test/test")
}

//
//
//
//// open .torrent file
//func Open(r io.Reader) (*bencodeTorrent, error) {
//	bto := bencodeTorrent{}
//	err := bencode.Unmarshal(r, &bto)
//	if err != nil {
//		return nil, err
//	}
//	return &bto, nil
//}
//
//func (bto *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
//
//}
//
//
//
///*解析peer信息*/
//func Unmarshal(peersBin []byte) ([]Peer, error) {
//	const peerSize = 6 // 4 字节是IP，2字节是Port
//	numPeers := len(peerBin) / peerSize
//	if len(peersBin) % peerSize != 0 {
//		err := fmt.Errorf("接受到的peer列表有错误！")
//        return nil, err
//	}
//	peers := make([]Peer, numPeers)
//	for i:=0; i < numPeers; i++ {
//		offset := i * peerSize
//		peers[i].IP = net.IP(peersBin[offset : offset+4])
//		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
//	}
//	return peers, nil
//}
//
//
//type Handshake struct {
//	Pstr 		string
//	InfoHash	[20]byte
//	PeerID		[20]byte
//}
//
//// 该成员主要为序列话Handshake
//func (h *Handshake) Serialize() []byte {
//	pstrlen := len(h.Pstr)
//	bufLen  := 49 + pstrlen  // 为什么49
//	buf     := make([]byte, bufLen)
//	buf[0]  := byte(pstrlen)
//	copy(buf[1:], h.Pstr)
//	// 空白8为
//	copy(buf[1+pstrlen+8:], h.InfoHash[:])
//	copy(buf[1+pstrlen+8+20:], h.PeerID[:])
//	return buf
//}
//
//// 从IO里面读数据然后反序列化
//func Read(r io.Reader) (*Handshake, error) {
//    // Do Serialize(), but backwards
//    // ...
//}
//
//
//type messageID uint8
//
///*消息类型的枚举*/
//const (
//	MsgChock			messageID = 0
//	MsgUnchoke       	messageID = 1
//    MsgInterested    	messageID = 2
//    MsgNotInterested 	messageID = 3
//    MsgHave          	messageID = 4
//    MsgBitfield      	messageID = 5
//    MsgRequest       	messageID = 6
//    MsgPiece         	messageID = 7
//    MsgCancel        	messageID = 8
//)
//
//type Message struct {
//	ID 		messageID
//	Payload	[]byte
//}
//
///* 序列化消息 */
//
//func (m *Message) Serialize() []byte {
//	if m == nil {
//		return make([]byte, 4) // 建立一个
//	}
//	length 	:= uint32(len(m.Payload) + 1)
//	buf 	:= make([]byte, 4+length)  // 长度 + ID + pylooad
//	binary.BigEndian.PutUint32(buf[0:4], length)
//	buf[4]  = byte(m.ID)
//	copy(buf[5:], m.Payload)
//	return buf
//}
//
//func Read(r io.Reader) (*Message, error) {
//	lengthBuf 	:= 	make([]byte, 4)  // 4B 的长度
//	_, err 		:= 	io.ReadFull(r, lengthBuf)  // 读取长度
//	if err != nil {
//		return nil, err
//	}
//	length 		:= binary.BigEndian.Uint32(lengthBuf)  // 转换为uint32
//
//	// keep-alive message
//	if length == 0 {
//		return nil, nil
//	}
//
//	messageBuf 	:= make([]byte, length)
//	_, err = io.ReadFull(r, messageBuf)
//	if err != nil {
//		return nil, err
//	}
//
//	// 返序列化
//	m := Message{
//		ID:			messageID(messageBuf[0])
//		Payload: 	messageBUf[1:]
//	}
//
//	return &m, nil
//}
//
//type Bitfield []byte
//
//func (bf Bitfield) HasPiece(index int) bool {
//	byteIndex 	:= index / 8    // 在那个byte存
//	offset 		:= index % 8	// 再偏离多少位
//	/*
//     *               *
//	 *         1 1 1 1   1 1 0 0
//	 *        ||<-o->|| ->||
//	 *                 7-offset
//	 */
//	return bf[byteIndex] >> (7 - offset) & 1 != 0
//}
//
//func (bf Bitfield) SetPiece(idx int) {
//	byteIndex		:= index / 8
//	offset 			:= index % 8
//	bf[byteIndex] 	|= 1 << (7-offset)
//}
//
//
///**Init queues for workers **/
//// 创建队列一个peer创建一个
//workQueue := make(chan *pieceWork, len(t.PieceHashes))
//
//results := make(chan *pieceResult)
//
//for idx, hash := range t.PieceHashes {
//	length := t.calculatePieceSize(idx) // 计算长度
//	workQueue <- &pieceWork{idx, hash, length} // 进队
//}
//
//// start workers
//for _, peer := range t.Peers {
//	go t.startDownloadWorker(peer, workQueue, results)
//}
//
//// 收集结果到一个buf
//
//buf := make([]byte, t.Length)
//donePieces := 0
//for donePieces < len(t.PieceHashes) {
//	res := <- results
//	begin, end = t.calculateBoundsForPiece(res.index)
//	copy(buf[begin:end], res.buf)
//	donePieces++
//}
//
//close(workQueue)
//
//
//func (t *Torrent) startDownloadWorker(peer peer.Peer, workQueue chan *pieceWork, results chan *pieceResult)  {
//	c, err	:= client.New(peer, t.PeerID, t.InfoHash)
//	if err != nil {
//		log.Printf("Could not handshake with %s, Disconnecting!\n", peer.IP)
//		return
//	}
//	defer c.Conn.Close()
//	log.Printf("Completed handshake with %s\n", peer.IP)
//
//	c.SendUnchoke()   	// 发送给peer unchoke 的标志
//	c.SendInterested()  // 发送活跃信息
//
//	for pw	:= range workQueue {
//		if !c.Bitfield.HasPiece(pw.index) { // 未完成
//			workQueue <- pw
//			continue
//		}
//	}
//	// 下载
//	buf, err := attemptDownloadPiece(c, pw)
//	if err != nil {
//		log.Println("Exiting", err)
//		workQueue <- pw // Put piece back on the queue
//		return
//	}
//	err = checkIntegrity(pw, buf)  // 检查
//	if err != nil {
//		log.Printf("Piece #%d failed integrity check\n", pw.index)
//		workQueue <- pw // Put piece back on the queue
//		continue
//	}
//
//	c.sendHave(pw.index)
//	results <- &pieceResult{pw.index, buf} // 结果进队
//
//}
//
///*
// *  1. 管理状态
// */
//
//type pieceProgress struct {
//	index		int
//	client		*client,Client
//	buf			[]byte
//	downloaded	int
//	requested	int
//	backlog		int
//}
//
//func (state *pieceProgress) readMessage() error {
//	msg, err := state.client.Read()
//	switch msg.ID {
//	case message.MsgUnchoke:           	// 标记
//		state.client.Choked = false
//	case message.MsgChoke:
//        state.client.Choked = true
//    case message.MsgHave:				// 查询
//        index, err := message.ParseHave(msg)
//        state.client.Bitfield.SetPiece(index)
//    case message.MsgPiece:				// Picese
//        n, err := message.ParsePiece(state.index, state.buf, msg)
//        state.downloaded += n
//        state.backlog--
//    }
//    return nil
//}
//
