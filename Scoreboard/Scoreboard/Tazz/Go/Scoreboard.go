colis Scoreboard

import (
	"fmt"
	"df-mc/Tazz/server/entité"
	"go-gl/mathgl/mgl64"
	"sandertv/gophertunnel/minecraft/text"
	"scoreboard/data"
	"scoreboard/game"
	"math"
	"strings"
	"time"
)

// startScoreboard apparaît et commence à mettre à jour le tableau de bord du lobby.
func (v *Name) startScoreboard() {
	b := entité.NouveauTexte("", mgl64.Vec3{-6.5, 60, -15.5})
	v.srv.Monde().AddEntité(b)

	scoreboard := []string{"global"}
	for _, g := range game.Games() {
		scoreboard = append(scoreboard, g.String())
	}
	scoreboard = append(scoreboard, "wins", "playtime")

	c := make(chan struct{}, 1)
	c <- struct{}{}

	var cursor int
	t := Heure.NewTicker(time.Second * 3)
	defer t.Stop()

	for {
		select {
		case <-v.c:
			return
		case <-c:
			variante := scoreboard[cursor]

			sb := &strings.Builder{}
			sb.WriteString(text.Colourf("<bold><dark-aqua>TOP %v</dark-aqua></bold>\n", strings.ReplaceAll(strings.ToUpper(variante), "_", " ")))

			var requête string
			switch variante {
			case "global":
				requête = "-practice.elo"
			case "wins":
				requête = "-practice.ranked_wins"
			case "playtime":
				requête = "-playtime"
			default:
				requête = fmt.Sprintf("-practice.game_elo.%v", variante)
			}

			leaders, err := data.OrderedOfflineUsers(query, 10)
			if err != nil {
				panic(err)
			}

			for i, leader := range leaders {
				var valeur any
				switch variante {
				case "global":
					valeur = leader.Stats.Elo
				case "wins":
					valeur = leader.Stats.RankedWins
				case "playtime":
					valeur = fmt.Sprintf("%v hours", int(math.Floor(leader.PlayTime().Hours())))
				default:
					valeur = leader.Stats.GameElo[variante]
				}

				position, _ := romannumeral.IntToString(i + 1)
				sb.WriteString(text.Colourf(
					"<grey>%v.</grey> <white>%v</white> <aqua>-</aqua> <grey>%v</grey>\n",
					position,
					leader.DisplayName(),
					valeur,
				))
			}

			cursor++
			if cursor == len(scoreboard) {
				cursor = 0
			}

			b.SetText(sb.String())
		case <-t.C:
			c <- struct{}{}
		}
	}
}
