package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path"
	"runtime"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func init() {
	app.Action = extract
	app.Flags = []cli.Flag{
		utils.DataDirFlag,
	}
}

func WriteTx(db *badger.DB, tx *types.Transaction) error {
	hash := tx.Hash()
	key := hash[:]
	value, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	txn := db.NewTransaction(true)
	if err := txn.Set(key, value); err != nil {
		return err
	}
	return txn.Commit()
}

func WriteTxs(db *badger.DB, height uint64, txs types.Transactions) error {
	hashes := make([]common.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.Hash()
	}
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, height)
	value, err := rlp.EncodeToBytes(hashes)
	if err != nil {
		return err
	}
	txn := db.NewTransaction(true)
	if err = txn.Set(key, value); err != nil {
		return err
	}
	if err = txn.Commit(); err != nil {
		return err
	}
	for _, tx := range txs {
		if err = WriteTx(db, tx); err != nil {
			return err
		}
	}
	return nil
}

const start, end uint64 = 19130000, 19140000

func extract(ctx *cli.Context) error {
	dataDir := ctx.String(utils.DataDirFlag.Name)
	db1, err := rawdb.Open(rawdb.OpenOptions{
		Directory:         path.Join(dataDir, "geth", "chaindata"),
		AncientsDirectory: path.Join(dataDir, "geth", "chaindata", "ancient"),
		ReadOnly:          true,
	})
	if err != nil {
		return err
	}
	defer db1.Close()

	blocks := make([]*types.Block, 0)
	for i := start; i < end; i++ {
		hash := rawdb.ReadCanonicalHash(db1, i)
		block := rawdb.ReadBlock(db1, hash, i)
		if math.Abs(float64(block.Size())-141.40*1024) < 1600 {
			blocks = append(blocks, block)
		}
	}

	rand.Shuffle(len(blocks), func(i, j int) {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	})

	if err = os.RemoveAll("txs"); err != nil {
		return err
	}
	db2, err := badger.Open(badger.DefaultOptions("txs").WithLoggingLevel(badger.ERROR))
	if err != nil {
		return err
	}
	defer db2.Close()
	size := uint64(0)
	for i := uint64(0); i < 100; i++ {
		block := blocks[i]
		if err = WriteTxs(db2, i, block.Transactions()); err != nil {
			return err
		}
		size += block.Size()
	}
	fmt.Println(float64(size) / 100 / 1024)

	return db2.Flatten(runtime.NumCPU())
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
