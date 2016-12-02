package candidates

import (
	"fmt"
	"strconv"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Candidates struct {
	Name string
	Votes int
}

/** 投票者マップ */
var candidatesMap = make(map[string]int)

func init() {
	http.HandleFunc("/vote", vote)
	http.HandleFunc("/getVotes", getVotes)
}

/**
 * vote
 * rest api
 * request param {name, num}
 */
func vote(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	// 投票者名を取得
	// TODO: 投票者のホワイトリスト処理を入れる
	name := r.URL.Query().Get("name")
	updateCandidatesByName(c, name)

	// 投票する票数を取得
	num, _ := strconv.Atoi(r.URL.Query().Get("num"))

	if (num < 0) {
		// マイナスはありえないのでそのまま返す
		writeResponse(c, w, name)
		return;
	} else if (1000 < num) {
		num = 1000
	}

	candidatesMap[name] = candidatesMap[name] + num

	// 投票処理
	var temp = Candidates {
		Name: name,
		Votes: candidatesMap[name],
	}
	// アクセスキーを取得
	key := datastore.NewKey(c, "Candidates", name, 0, nil)
	datastore.Put(c, key, &temp)

	// レスポンスを返す
	writeResponse(c, w, name)
}

/**
 * getVotes
 * rest api
 * request param {name}
 * return {"response": {"name": name, "votes": votes}}
 */
func getVotes(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// 投票者名を取得
	// TODO: 投票者のホワイトリスト処理を入れる
	name := r.URL.Query().Get("name")
	updateCandidatesByName(c, name)

	// レスポンスを返す
	writeResponse(c, w, name);
}

func updateCandidatesByName(c context.Context, name string) {
	if (name == "") {
		return
	}

	// 投票者のキャッシュマップに存在するか確認
	_, ok := candidatesMap[name]
	if (!ok) {
		// datastoreから問い合わせる
		var temp Candidates
		// アクセスキーを取得
		key := datastore.NewKey(c, "Candidates", name, 0, nil)
		error := datastore.Get(c, key, &temp)

		// datastoreにまだ値が無い場合
		if (error == nil) {
			candidatesMap[name] = temp.Votes
		} else {
			candidatesMap[name] = 0
		}
	}
}

func writeResponse(c context.Context, w http.ResponseWriter, name string) {
	// レスポンスを返す
	updateCandidatesByName(c, name)
	fmt.Fprint(w, "{\"response\": {\"name\": \"" + name + "\", \"votes\": " + strconv.Itoa(candidatesMap[name]) + "}}")
}