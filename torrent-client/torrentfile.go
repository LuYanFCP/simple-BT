package BTClient

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const HASH_LENGTH int = 20 // 20B
const SERVER_PORT uint16 = 50001 //

type TorrentFile struct {
	/* TorrentFile的对象 */
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

/* bencodeInfo的对象 */
type bencodeInfo struct {
	/* info 中解析出来的
	['files', 'name', 'piece length', 'pieces', 'private', 'source']
	*/
	Pieces      string `pieces`
	PieceLength int    `piece length` // 每个文件块的大小，用Byte计算
	Length      int    `length`       //
	Name        string `name`         //  推荐文件名
}

/* Torrent的从Torrentfile解析bencodeInfo */
type bencodeTorrent struct {
	/* 一般来说种子第一步解析出来的文件
	['ED82F155', 'announce', 'created by', 'created rd', 'created type', 'creation date', 'encoding', 'info']
	*/
	Announce string      `announce`
	Info     bencodeInfo `info`
}

/*
 *  resp
 *
 */

type bencodeTrackerResp struct {
	Interval int    `interval`
	Peer     string `peers`
}


/*#####################################################
 *#          TorrentFile文件解析、下载                  #
 *#####################################################
 */

/* 打开并解析种子文件 */
func OpenTorrentFile(path string) (TorrentFile, error) {
	reader, err := os.Open(path) // 打开
	if err != nil {
		log.Fatalln("打开文件失败，请查看路径是否正确！", err)
		return TorrentFile{}, err
	}
	defer reader.Close() // 关掉fptr， 这里考虑一下

	// 解析相应的bencodeTorrent
	btf := bencodeTorrent{}
	err = bencode.Unmarshal(reader, &btf)
	if err != nil {
		log.Fatalln("bencode解析错误", err)
		return TorrentFile{}, nil
	}
	return btf.parse() // 解析
}

func (trf *TorrentFile) Download(path string) error {
	var peerID [20]byte
	copy(peerID[:8], []byte("-lt0D80-"))  // 伪造一个rtorrent的开头

	/*
	 *  BT:https://blog.csdn.net/pi408637535/article/details/44731393
	 */
	_, err := rand.Read(peerID[8:])  // 后12位随机
	// 定死 peerID， BT要求一个peerID只能请求3次
	peerID = [20]byte{45, 108, 116, 48, 68, 56, 48, 45, 128, 49, 24, 54, 165, 225, 41, 66, 25, 184, 153, 164}
	fmt.Println(peerID)
	if err != nil {
		return err
	}
	// 从服务器get Peers的列表
	peers, err := trf.getPeers(peerID, SERVER_PORT)
	if err != nil {
		return err
	}
	// 打印
	for i:=0; i < len(peers); i++ {
		fmt.Println(string(peers[i].IP), peers[i].Port)
	}

	return err
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	/*
		成员函数， 构建请求的url
	*/
	base, err := url.Parse(t.Announce) // 构建基础url
	if err != nil {
		return "", err
	}
	/* 构建请求tracker的信息 **/
	real_hash, _:= hex.DecodeString("e602b4193f4c099fc9164381de01b790def7facf")
	fmt.Println(real_hash)
	fmt.Println(t.InfoHash)
	//copy(t.InfoHash, real_hash)
	//string(peerID[:])
	params := url.Values{
		"info_hash":  []string{string(real_hash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"numwant":    []string{"80"},
		"compact":    []string{"1"},
		"left":       []string{"0"}, //strconv.Itoa(t.Length)
	}
	base.RawQuery +=  ("&" + params.Encode())
	return base.String(), nil // 构建请求的url
}

func peerUnmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 18   // 4个ip， 2个port
	numsPeers := len(peersBin) / peerSize
	if len(peersBin) % peerSize != 0 {

	}
	peers := make([]Peer, numsPeers)
	for i:=0; i < numsPeers; i++ {
		offset := i * peerSize
		peers[i].IP   = net.IP(peersBin[offset: offset+16])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBin[offset+16:offset+18]))
	}
	return peers, nil
}

func (trf *TorrentFile) getPeers(peerID [20]byte, port uint16) ([]Peer, error) {
	url, err := trf.buildTrackerURL(peerID, port) // 获取到url
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "rtorrent/0.9.8/0.13.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Acept", "*/*")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("超时异常")
		return nil, nil
	}
	defer resp.Body.Close()

	respBody := resp.Body

	// 有可能使用gzip encode
	if resp.Header.Get("Content-Encoding") == "gzip" {
		respBody, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal("gzip解析错误!")
			return nil, nil
		}
	}

	trackerResp := bencodeTrackerResp{}

	err = bencode.Unmarshal(respBody, &trackerResp)
	if err != nil {
		return nil, err
	}
	return peerUnmarshal([]byte (trackerResp.Peer))
}

/* ##########################################################
 * #              bemcodeTorrent的成员函数                   #
 * ##########################################################
 */

/* 解析bencode生成TorrentFile的文件 */
func (bct *bencodeTorrent) parse() (TorrentFile, error) {
	/*
		Announce    string
		InfoHash    [20]byte
		PieceHashes [][20]byte
		PieceLength int
		Length      int
		Name        string
	*/
	infoHash, err := bct.Info.hash() // 计算hash
	if err != nil {
		log.Fatalln("hashe 解析的问题！", err)
		return TorrentFile{}, err
	}
	pieceHashes, err := bct.Info.getPieceHashes()
	if err != nil {
		log.Fatalln("获得Piceses hash解析出现了问题！", err)
		return TorrentFile{}, err
	}

	tf := TorrentFile{
		Announce:    bct.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bct.Info.PieceLength,
		Length:      bct.Info.Length,
		Name:        bct.Info.Name,
	}
	return tf, nil
}

/* ##########################################################
 * # Errorf             bencodeInfo的成员函数                #
 * ##########################################################
 * 提供的功能：
 * 		1. 获取获取整体的hash值， 主要构成 TorrentFile文件
 *      2. 分割pieces中的hash值， 分割成列表， 这里使用的hash函数是sha1->20
 *
 */

func (bi *bencodeInfo) hash() ([20]byte, error) {
	// err := bencode.Marshal(writer, data)
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *bi) // 编码，然后hash
	if err != nil {
		log.Fatal("bencode编码的时候出现了问题！", err)
		return [20]byte{}, err
	}
	sha1Val := sha1.Sum(buf.Bytes())
	return sha1Val, nil
}

func (bi *bencodeInfo) getPieceHashes() ([][20]byte, error) {
	buf := []byte(bi.Pieces) // 没有拆分开的pieces
	if len(buf)%HASH_LENGTH != 0 {
		log.Fatalln("bt文件的格式不对，hash值对不上!")
		err := fmt.Errorf("bt文件的格式不对，hash值对不上! 长度:%d", len(buf))
		return [][20]byte{}, err
	}
	/* 进行切割 */
	numsPieces := len(bi.Pieces) / HASH_LENGTH
	hashes := make([][20]byte, numsPieces) // 长numsPieces的数组
	offset := 0
	for i := 0; i < numsPieces; i++ {
		offset = i * HASH_LENGTH
		copy(hashes[i][:], buf[offset:offset+HASH_LENGTH])
	}
	return hashes, nil
}
