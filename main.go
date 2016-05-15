package main

import (
    "fmt"
    "log"
    "time"
    "math/rand"
    "strings"
    "os"
    "strconv"
)

type Player struct {
    money byte
    bids []byte
    drawAdv bool
}
func (me *Player) toString(moveNum int) string {
    adv := ""
    if me.drawAdv {
        adv = "*"
    }

    return fmt.Sprintf("% 4d % 4d % 4s", me.bids[moveNum], me.money, adv)
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    fmt.Println("The Bidding Game tester.")

    numAttempts := 0
    if len(os.Args) >= 2 {
        var err error
        numAttempts, err = strconv.Atoi(os.Args[1])
        if err != nil {
            log.Fatalln("Invalid argument - must be integer")
        }
    }

    fmt.Println("Attempts:", numAttempts)
    fmt.Println("")

    winnersFile, ferr := os.OpenFile("winners.csv", os.O_CREATE | os.O_APPEND | os.O_WRONLY, os.ModePerm)
    if ferr != nil {
        log.Fatal(ferr)
    }
    defer winnersFile.Close()

    for attempt := 0; attempt < numAttempts; attempt++ {
        drawAdvRand := rand.Intn(2)
        winner, err := playGame(100, drawAdvRand)
        if err != nil {
            log.Print(err)
        }
        if winner != nil {
            winnerLine := strings.Replace(fmt.Sprintf("%v", winner.bids), " ", ",", -1)
            winnersFile.WriteString(winnerLine[1:len(winnerLine)-1] + "\n")
        }
    }
}

func playGame(initMoney byte, drawAdvantage int) (*Player, error) {
    bottlePos := 5

    pL := Player{initMoney, make([]byte, 1), drawAdvantage == 0}
    pR := Player{initMoney, make([]byte, 1), drawAdvantage == 1}

    fmt.Println("Move         Left           Right     L-1-2-3-4-5-6-7-8-9-R")
    fmt.Println("       Bid    $   DA   Bid    $   DA")

    moveNum := 0
    for {
        bottlePosImg := strings.Repeat("-", bottlePos*2) + fmt.Sprintf("%v", bottlePos) + strings.Repeat("-", 20 - bottlePos*2)
        fmt.Printf("% 4d  %s  %s  %s\n",
            moveNum,
            pL.toString(moveNum),
            pR.toString(moveNum),
            bottlePosImg)

        if bottlePos == 0 || bottlePos == 10 || pL.money == 0 && pR.money == 0 {
            fmt.Println(strings.Repeat("=", 60))

            if bottlePos == 0 {
                fmt.Println("Player L wins!")
                return &pL, nil
            } else if bottlePos == 10 {
                fmt.Println("Player R wins!")
                return &pR, nil
            } else if pL.money == 0 && pR.money == 0 {
                fmt.Println("Draw!")
                return nil, nil
            } else {
                log.Fatalln("Invalid win state", bottlePos, pL.money, pR.money)
            }
        }

        if (moveNum > 1000) {
            return nil, fmt.Errorf("Move number = %d", moveNum)
        }

        bottlePos = move(&pL, &pR, bottlePos)
        moveNum = moveNum+1
    }

    // return nil, fmt.Errorf("Invalid game state")
}

func doBid(p *Player) byte {
    if p.money <= 0 {
        return 0
    }

    r := rand.Intn(int(p.money)) + 1 // [0, p.money) -> [1, p.money+1) === [1, p.money]

    return byte(r)
}

func move(p1 *Player, p2 *Player, bottlePos int) int {
    b1 := doBid(p1)
    b2 := doBid(p2)

    p1.bids = append(p1.bids, b1)
    p2.bids = append(p2.bids, b2)

    p1winner := func() int {
        p1.money -= b1
        return bottlePos-1
    }
    p2winner := func() int {
        p2.money -= b2
        return bottlePos+1
    }

    if b1 > b2 {
        return p1winner()
    } else if b1 < b2 {
        return p2winner()
    } else {
        if p1.drawAdv && !p2.drawAdv {
            p1.drawAdv = false
            p2.drawAdv = true
            return p1winner()
        } else if p2.drawAdv && !p1.drawAdv {
            p2.drawAdv = false
            p1.drawAdv = true
            return p2winner()
        } else {
            log.Fatalf("Invalid game state: p1.drawAdv=%v, p2.drawAdv=%v\n", p1.drawAdv, p2.drawAdv)
        }
    }

    return 0
}
