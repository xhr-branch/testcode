package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

// 区块结构
type Block struct {
	Index        int
	PreviousHash string
	Timestamp    string
	Data         string
	Hash         string
	Nonce        int
}

// 区块链结构
type Blockchain struct {
	Blocks []*Block
}

// 节点结构
type Node struct {
	ID         string
	Reputation int
	Authority  bool
}

// 计算 SHA-256 哈希
func calculateHash(block *Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.PreviousHash, block.Timestamp, block.Data, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

// 创建创世区块
func createGenesisBlock() *Block {
	return &Block{
		Index:        0,
		PreviousHash: "",
		Timestamp:    time.Now().String(),
		Data:         "Genesis Block",
		Hash:         "",
		Nonce:        0,
	}
}

// 创建新区块
func createBlock(previousBlock *Block, data string) *Block {
	block := &Block{
		Index:        previousBlock.Index + 1,
		PreviousHash: previousBlock.Hash,
		Timestamp:    time.Now().String(),
		Data:         data,
		Hash:         "",
		Nonce:        0,
	}
	block.Hash = calculateHash(block)
	return block
}

// 挖矿（找到符合条件的哈希值）
func mineBlock(block *Block, difficulty int) {
	target := strings.Repeat("0", difficulty) // 目标是指定前面为0的哈希
	for {
		block.Hash = calculateHash(block)
		if block.Hash[:difficulty] == target {
			break
		}
		block.Nonce++
	}
}

// 创建区块链
func createBlockchain() *Blockchain {
	blockchain := &Blockchain{}
	genesisBlock := createGenesisBlock()
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
	return blockchain
}

// 添加区块到区块链
func (bc *Blockchain) addBlock(newBlock *Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

// 保存区块链到文件
func saveBlockchainToFile(blockchain *Blockchain, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(blockchain)
	if err != nil {
		return err
	}
	return nil
}

// 从文件读取区块链
func loadBlockchainFromFile(filename string) (*Blockchain, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var blockchain Blockchain
	err = decoder.Decode(&blockchain)
	if err != nil {
		return nil, err
	}
	return &blockchain, nil
}

// 创建 Web 服务器，处理请求并展示区块链数据
func serveBlockchain(w http.ResponseWriter, r *http.Request) {
	// 尝试从文件加载区块链
	filename := "blockchain.dat"
	blockchain, err := loadBlockchainFromFile(filename)
	if err != nil {
		// 如果文件加载失败，则创建一个新的区块链
		blockchain = createBlockchain()
	}

	difficulty := 4

	// 假设有5个节点，包含不同的声誉评分和权威性
	nodes := []Node{
		{"Node1", 50, true},
		{"Node2", 30, false},
		{"Node3", 80, true},
		{"Node4", 60, false},
		{"Node5", 90, true},
	}

	// 挖矿并添加新区块
	for i := 1; i <= 5; i++ {
		// 选择权威节点
		selectedNode := nodes[i%5] // 模拟节点选择

		// 创建新区块
		newBlock := createBlock(blockchain.Blocks[len(blockchain.Blocks)-1], fmt.Sprintf("Block %d Data", i))

		// 挖矿：找到符合条件的哈希值
		mineBlock(newBlock, difficulty)

		// 将新区块添加到区块链
		blockchain.addBlock(newBlock)

		// 输出矿工信息
		fmt.Printf("Block %d mined by %s with Reputation %d\n", newBlock.Index, selectedNode.ID, selectedNode.Reputation)
	}

	// 保存区块链到文件
	err = saveBlockchainToFile(blockchain, filename)
	if err != nil {
		http.Error(w, "Failed to save blockchain", http.StatusInternalServerError)
		return
	}

	// 创建页面数据
	type PageData struct {
		Blockchain *Blockchain
		Nodes      []Node // 添加节点信息
	}

	// 打印区块链的 HTML 模板
	const htmlTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Blockchain</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			.block { margin-bottom: 20px; padding: 15px; border: 1px solid #ccc; }
			.block h3 { margin: 0; }
			.block p { margin: 5px 0; }
			.node { margin-bottom: 15px; padding: 15px; border: 1px solid #ddd; }
		</style>
	</head>
	<body>
		<h1>Blockchain</h1>
		<h2>Nodes</h2>
		{{range .Nodes}}
		<div class="node">
			<h3>{{.ID}}</h3>
			<p><strong>Reputation:</strong> {{.Reputation}}</p>
			<p><strong>Authority:</strong> {{.Authority}}</p>
		</div>
		{{end}}

		<h2>Blockchain</h2>
		{{range .Blockchain.Blocks}}
		<div class="block">
			<h3>Block #{{.Index}}</h3>
			<p><strong>Timestamp:</strong> {{.Timestamp}}</p>
			<p><strong>Data:</strong> {{.Data}}</p>
			<p><strong>Hash:</strong> {{.Hash}}</p>
			<p><strong>Previous Hash:</strong> {{.PreviousHash}}</p>
			<p><strong>Nonce:</strong> {{.Nonce}}</p>
		</div>
		{{end}}
	</body>
	</html>
	`

	// 渲染 HTML
	tmpl, err := template.New("blockchain").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Blockchain: blockchain,
		Nodes:      nodes, // 传递节点信息
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// 启动 Web 服务器
	http.HandleFunc("/", serveBlockchain)
	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

