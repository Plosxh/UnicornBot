package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	//"go/types"
	"os/exec"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Variables used for command line parameters
var (
	Token              string
	BotID              string
	Guild_Roles        []*discordgo.Role
	MonstersList		map[int][]string
	EmbeeFilter			bool
	RunesMapping		[]RMapping
	AttributeMapping		[]string
)

type Nstring string
type RMapping struct {
	Type string
	Set int
}

type errorString struct {
	s string
}

type SwarfarmResponseProfile struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []ProfilMonster `json:"results"`
}

type Profile struct {
	Monster int  `json:"monster"`
	Stars int `json:"stars"`
	Level int `json:"level"`
}

type SwarfarmResponse struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []Monster `json:"results"`
}
type Monster struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Image string `json:"image_filename"`
	Element string `json:"element"`
	Archetype string `json:"archetype"`
	NaturalStars int `json:"natural_stars"`
	SkillUpsToMax int `json:"skill_ups_to_max"`
	LeaderSkill LeaderSkill `json:"leader_skill"`
	Skills []int `json:"skills"`
	HP int `json:"base_hp"`
	ATK int `json:"base_attack"`
	DEF int `json:"base_defense"`
	SPD int `json:"speed"`
	CR int `json:"crit_rate"`
	CD int `json:"crit_damage"`
	RES int `json:"resistance"`
	ACC int `json:"accuracy"`

}

type Rune struct {
	ID string  `json:"id"`
	Type int `json:"type"`
	Slot int `json:"slot"`
	Stars int `json:"stars"`
	Level int `json:"level"`
	Quality int `json:"quality"`
	OriginalQuality int `json:"original_quality"`
	Value int `json:"value"`
	SubstatUpgradeRemaining int `json:"substat_upgrades_remaining"`
	Efficiency float64 `json:"efficiency"`
	EfficiencyMax float64 `json:"max_efficiency"`
	Main int `json:"main_stat"`
	MainValue int `json:"main_stat_value"`
	Innate int `json:"innate_stat"`
	InnateValue int `json:"innate_stat_value"`
	Substat1 int `json:"substat_1"`
	Substat1Value int `json:"substat_1_value"`
	Substat1Craft int `json:"substat_1_craft"`
	Substat2 int `json:"substat_2"`
	Substat2Value int `json:"substat_2_value"`
	Substat2Craft int `json:"substat_2_craft"`
	Substat3 int `json:"substat_3"`
	Substat3Value int `json:"substat_3_value"`
	Substat3Craft int `json:"substat_3_craft"`
	Substat4 int `json:"substat_4"`
	Substat4Value int `json:"substat_4_value"`
	Substat4Craft int `json:"substat_4_craft"`

}

type Building struct {
	Id int `json:"id"`
	Area string `json:"area"`
	AffectedStat string `json:"affected_stat"`
	Element string `json:"element"`
	Name string `json:"name"`
	MaxLevel int `json:"max_level"`
	StatBonus []int `json:"stat_bonus"`
	UpgradeCost []int `json:"upgrade_cost"`
	Description string `json:"description"`
}

type Buildings struct {
	Count int `json"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []Building `json:"results"`
}

type swarfarmBuilding struct {
	Id string `json:"id"`
	URL string `json:"url"`
	BuildingID int `json:"building"`
	Level int `json:"level"`
}

type ProfilBuilding struct {
	Count int `json:"int"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []swarfarmBuilding `json:"results"`
}

type ProfilMonster struct {
	Id int `json:"monster"`
	Name string `json:"name"`
	Level int `json:"level"`
	Star int `json:"stars"`
	//Image string `json:"image_filename"`
	Skill1Level int `json:"skill_1_level"`
	Skill2Level int `json:"skill_2_level"`
	Skill3Level int `json:"skill_3_level"`
	Skill4Level int `json:"skill_4_level"`
	SkillUpsToMax int `json:"skill_ups_to_max"`
	BaseHP int `json:"base_hp"`
	BaseATK int `json:"base_attack"`
	BaseDEF int `json:"base_defense"`
	BaseSPD int `json:"base_speed"`
	BaseCR int `json:"base_crit_rate"`
	BaseCD int `json:"base_crit_damage"`
	BaseRES int `json:"base_resistance"`
	BaseACC int `json:"base_accuracy"`
	RuneHP int `json:"rune_hp"`
	RuneATK int `json:"rune_attack"`
	RuneDEF int `json:"rune_defense"`
	RuneSPD int `json:"rune_speed"`
	RuneCR int `json:"rune_crit_rate"`
	RuneCD int `json:"rune_crit_damage"`
	RuneRES int `json:"rune_resistance"`
	RuneACC int `json:"rune_accuracy"`
	AverageRuneffeciency float64 `json:"avg_rune_efficiency"`
	Runes []Rune `json:"runes"`
}

type LeaderSkill struct {
	Id int `json:"id"`
	Attribute string `json:"attribute"`
	Amount int `json:"amount"`
	Area string `json:"area"`
	Element string `json:"element"`
}

type Skill struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CoolTime int `json:"cooltime"`
	Hits int `json:"hits"`
	MaxLevel int `json:"max_level"`
	Formula string `json:"multiplier_formula"`
	Skills []int `json:"skill"`
	Image string `json:"icon_filename"`
	SkillUp []string `json:"level_progress_description"`
}

// Todo
// Buildings
// add skill up dans ?skill
// Teams


func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
}

func main() {
	MonstersList = make(map[int][]string)
	EmbeeFilter = false
	AttributeMapping = []string{"HP","HP%","ATK","ATK%","DEF","DEF%","SPD","CR","CD","RES","ACC"}
	//MonstersList[970] = "Adrian"
	//MonstersList[1] = "Forest Keeper"
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + "NjE4NDE3NjcwNTIxMzU2MzA4.XW5Yfw.P3QtTCg2IyioWL8T0J4mx7zSR0g")
	//dg, err := discordgo.New("Bot " + "NjA2NDUxNTMyNDE3OTkwNjU2.XZRhLw.Fg7qhPaQ3fjCdYZEhVLMlRNBFl4")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details try again,", err)
	}
	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(UseReaction)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.Contains(m.Content, "?") {
		// Ignore all messages created by the bot itself
		// chann, errorChan := s.Channel(m.ChannelID)
		// if errorChan != nil {
		// 	fmt.Println(errorChan)
		// }
		if m.Author.ID == "400927402106028034" && EmbeeFilter{
			_, _ = s.ChannelMessageSend(m.ChannelID, "please stop Bothering me ! :angry:")
		}else if m.Author.ID == "210801951552700416" && strings.Contains(m.Content,"?Embee"){
			if EmbeeFilter {
				EmbeeFilter = false
				_, _ = s.ChannelMessageSend(m.ChannelID, "<@400927402106028034> has been freed")
			}else{
				EmbeeFilter = true
				_, _ = s.ChannelMessageSend(m.ChannelID, "<@400927402106028034> has been blocked")
			}
		}else if strings.Contains(m.Content,"?update") && m.Author.ID == "210801951552700416" {
			out, err := exec.Command("git", "pull").Output()
			if err != nil {
				print(err.Error())
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, string(out))

		}else if strings.Contains(m.Content,"?add") && m.Author.ID == "210801951552700416" {
			var GuildeMembers []string
			var profile SwarfarmResponseProfile
			splittedString := strings.TrimPrefix(m.Content, "?add ")
			jsonFile, err := os.Open("guilde.json")
			if err != nil {
				fmt.Println(err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &GuildeMembers)
			errJson := getJson("https://swarfarm.com"+"/api/v2/profiles/"+splittedString, &profile)
			if errJson != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, splittedString+" profile does not exist")
			}else if !stringInSlice(splittedString, GuildeMembers){
				GuildeMembers = append(GuildeMembers, splittedString)
				fmt.Println(GuildeMembers)
				file,errorJson := json.Marshal(GuildeMembers)
				erroFile := ioutil.WriteFile("./guilde.json",file,0644)
				if erroFile != nil{
					fmt.Println(erroFile)
				}
				if errorJson != nil{
					fmt.Println(errorJson)
				}
				_, _ = s.ChannelMessageSend(m.ChannelID, splittedString+" successfuly added !")
			}else{
				_, _ = s.ChannelMessageSend(m.ChannelID, splittedString+" already in base !")
			}
		}else if strings.Contains(m.Content,"?monster"){
			var monsters Monster
			jsonFile, err := os.Open("monsters.json")
			// if we os.Open returns an error then handle it
			if err != nil {
				fmt.Println(err)
			}
			// defer the closing of our jsonFile so that we can parse it later on
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &MonstersList)
			//splittedString := strings.Split(m.Content, " ")
			url := "https://swarfarm.com/api/v2/monsters/"
			splittedString := strings.TrimPrefix(m.Content, "?monster ")
			for key, value := range MonstersList {
				if value[0] == splittedString {
					errJson := getJson(url+strconv.Itoa(key), &monsters)
					if errJson != nil {
						fmt.Println(errJson)
					}
					fields := make(map[string]string)
					//for key, value := range monsters.Skills {
					//	skill, errSkill := GetSkill(value)
					//	if errSkill != nil{
					//		fmt.Println(errSkill)
					//	}else {
					//		fields["Skill "+strconv.Itoa(key+1)] = skill.Print()
					//	}
					//}
					fields["Principales"] ="**HP:** "+strconv.Itoa(monsters.HP)+
						"\n **ATK:** "+ strconv.Itoa(monsters.ATK)+
						"\n **DEF:** "+strconv.Itoa(monsters.DEF)+
						"\n **SPD:** "+strconv.Itoa(monsters.SPD)
					fields["Secondaires"] ="**CR:** "+strconv.Itoa(monsters.CR)+
						"\n **CD:** "+strconv.Itoa(monsters.CD)+
						"\n **RES:** "+strconv.Itoa(monsters.RES)+
						"\n **ACC:** "+strconv.Itoa(monsters.ACC)

					ToSend := CreateAnEmbed(" Nat "+strconv.Itoa(monsters.NaturalStars),fields,monsters.Name,MonstersList[monsters.Id][1], "monsters")
					_, _ = s.ChannelMessageSendEmbed(m.ChannelID, ToSend)
					break
				}
			}

			//for i:=0;i<len(monsters.Results) ;i++  {
			//	if monsters.Results[i].Name == splittedString {
			//		_, _ = s.ChannelMessageSend(m.ChannelID, monsters.Results[i].PrintMonster())
			//	}else{
			//		//_, _ = s.ChannelMessageSend(m.ChannelID, "Je n'ai pas trouvé "+splittedString[1])
			//	}
			//}

		}else if strings.Contains(m.Content, "?gibsette"){
			gif := "https://tenor.com/view/ncis-slap-silly-scold-gibbs-gif-13523518"
			_, _ = s.ChannelMessageSend(m.ChannelID, gif)
		} else if strings.Contains(m.Content, "?sf"){
			///api/v2/profiles/{user_pk}/monsters/
			//splittedString := strings.Split(m.Content, " ")
			splittedString := strings.TrimPrefix(m.Content, "?sf ")
			userName := strings.Split(splittedString," ")
			if len(userName) >= 2{
				monstersHit := GetMonsters(url.QueryEscape(strings.TrimPrefix(splittedString, userName[0]+" ")), userName[0])
				IsMultiple := false
				if len(monstersHit) == 0 {
					_, _ = s.ChannelMessageSend(m.ChannelID, userName[0]+ " has no "+ url.QueryEscape(strings.TrimPrefix(splittedString, userName[0]+" ")))
				}else {
					fields, _, stars, level := monstersHit[0].FormatPrint()
					ToAdd := ""
					if len(monstersHit) > 1 {
						ToAdd = " 1/"+strconv.Itoa(len(monstersHit))
						IsMultiple = true
					}
					ToSend := CreateAnEmbed("Monster "+url.QueryEscape(strings.TrimPrefix(splittedString, userName[0]+" "))+ ToAdd +" Lvl "+level+" "+stars+"*",fields,userName[0],MonstersList[monstersHit[0].Id][1],"monsters")
					message, _ := s.ChannelMessageSendEmbed(m.ChannelID, ToSend)
					if IsMultiple {
						errorAddEmoji := s.MessageReactionAdd(m.ChannelID, message.ID,"⬅")
						errorAddEmoji = s.MessageReactionAdd(m.ChannelID, message.ID,"➡")
						if errorAddEmoji != nil {
							fmt.Println(errorAddEmoji)
						}
					}

				}
			}
		}else if strings.Contains(m.Content, "?building"){
			splittedString := strings.TrimPrefix(m.Content, "?building ")
			userName := strings.Split(splittedString," ")
			buildings,errorBuild := getBuildings(userName[0])
			if errorBuild != nil {
				fmt.Println(errorBuild)
			}
			fields, _, _ := buildings.FormatPrint()
			ToSend := CreateAnEmbed("Buildings",fields,userName[0],"","monsters")
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, ToSend)

		}else if strings.Contains(m.Content, "skill"){
			//_, _ = s.ChannelMessageSend(m.ChannelID, "Command under maintenance")
			///api/v2/profiles/{user_pk}/monsters/
			//splittedString := strings.Split(m.Content, " ")
			if strings.Contains(m.Content,"1") {
				splittedString := strings.TrimPrefix(m.Content, "?skill1 ")
				monster , errorMonster :=GetMonster(splittedString)
				if errorMonster != nil{
					_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
				}else{
					fmt.Println(monster.LeaderSkill)
					skill, errSkill := GetSkill(monster.Skills[0])
					if errSkill != nil{
						fmt.Println(monster.Skills[0])
						fmt.Println(errSkill)
					}
					_, _ = s.ChannelMessageSendEmbed(m.ChannelID, skill.Print())
				}


			}else if strings.Contains(m.Content,"2") {
				splittedString := strings.TrimPrefix(m.Content, "?skill2 ")
				monster , errorMonster :=GetMonster(splittedString)
				if errorMonster != nil{
					_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
				}else{
					skill, errSkill := GetSkill(monster.Skills[1])
					if errSkill != nil{
						fmt.Println(errSkill)
					}
					_, _ = s.ChannelMessageSendEmbed(m.ChannelID, skill.Print())
				}

			}else if strings.Contains(m.Content,"3") {
				splittedString := strings.TrimPrefix(m.Content, "?skill3 ")
				monster , errorMonster :=GetMonster(splittedString)
				if errorMonster != nil{
					_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
				}else{
					skill, errSkill := GetSkill(monster.Skills[2])
					fmt.Println(skill)
					if errSkill != nil{
						_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
					}
					fmt.Println("will print skill")
					toto, tutu := s.ChannelMessageSendEmbed(m.ChannelID, skill.Print())
					fmt.Println(toto)
					fmt.Println(tutu)
				}

			}else if strings.Contains(m.Content,"4") {
				splittedString := strings.TrimPrefix(m.Content, "?skill4 ")
				monster , errorMonster :=GetMonster(splittedString)
				if errorMonster != nil{
					_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
				}else{
					fmt.Println(len(monster.Skills))
					skill, errSkill := GetSkill(monster.Skills[3])
					if errSkill != nil{
						fmt.Println(errSkill)
					}
					_, _ = s.ChannelMessageSendEmbed(m.ChannelID, skill.Print())
				}

			}else if strings.Contains(m.Content,"leader") {
				splittedString := strings.TrimPrefix(m.Content, "?leaderskill ")
				monster , errorMonster :=GetMonster(splittedString)
				if errorMonster != nil{
					_, _ = s.ChannelMessageSend(m.ChannelID, errorMonster.Error())
				}else{
					var newLeader LeaderSkill
					if monster.LeaderSkill == newLeader{
						_, _ = s.ChannelMessageSend(m.ChannelID, splittedString+" has no leader skill")

					}else{
						fmt.Println(monster.LeaderSkill)
						_, _ = s.ChannelMessageSend(m.ChannelID, monster.LeaderSkill.Print())
					}
				}

			}
		}
	}

}

func GetMonsters(name, user string) []ProfilMonster {
	var monstersHit []ProfilMonster
	var profile SwarfarmResponseProfile
	url := "https://swarfarm.com"+"/api/v2/profiles/"+user+"/monsters/?monster__name="+name+"&ordering=avg_rune_efficiency"
	errJson := getJson(url,&profile)
	if errJson != nil {
		fmt.Println(errJson)
	}
	//monsterName := strings.TrimPrefix(splittedString, userName[0]+" ")

	if len(MonstersList) == 0{
		GetMonsterFile()
	}
	monsterId := GetMonsterId(name)
	//var monstersHit []ProfilMonster
	for i:=0;i<len(profile.Results) ;i++  {

		if profile.Results[i].Id == monsterId {
			fmt.Println("found a monster")
			monstersHit = append(monstersHit,profile.Results[i])
			fmt.Println(monstersHit)

		}
	}
	OrderedMonsters := reverseInts(monstersHit)
	return OrderedMonsters
}

func UseReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd)  {
	if m.UserID != BotID {
		message, error := s.ChannelMessage(m.ChannelID,m.MessageID)
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println("Ò")
		fmt.Println(m.Emoji)
		fmt.Println(m.Emoji.ID)
		if message.Author.ID == BotID {
			for i:= 0; i<len(message.Embeds) ; i++ {
				InitialMobFound := 0
				splittedString := strings.Split(message.Embeds[i].Description, " ")
				monsterName := url.QueryEscape(splittedString[1])
				url := "https://swarfarm.com"+"/api/v2/profiles/"+message.Embeds[i].Title+"/monsters/?monster__name="+monsterName
				fmt.Println(url)
				for s := 0; s<len(splittedString) ; s++ {
					if strings.Contains(splittedString[s], "/") {
						InitialMobFound,_ = strconv.Atoi(strings.Split(splittedString[s], "/")[0])
					}
				}
				monstersHit := GetMonsters(monsterName,message.Embeds[i].Title)
				if len(monstersHit) > InitialMobFound {
					fields, _, stars, level := monstersHit[InitialMobFound +1].FormatPrint()
					ToAdd := ""
					if len(monstersHit) > 1 {
						ToAdd = " "+strconv.Itoa(InitialMobFound +1)+"/"+strconv.Itoa(len(monstersHit) -1)
					}
					ToSend := CreateAnEmbed("Monster "+monsterName+ ToAdd +" Lvl "+level+" "+stars+"*",fields,message.Embeds[i].Title,MonstersList[monstersHit[0].Id][1],"monsters")
					//_, _ = s.ChannelMessageSendEmbed(m.ChannelID, ToSend)
					_,_ = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, ToSend)
				}

			}
			fmt.Println(m.Emoji)
			errorRemove := s.MessageReactionRemove(m.ChannelID,m.MessageID,m.Emoji.Name,m.UserID)
			if errorRemove != nil {
				fmt.Println(errorRemove)
			}
		}
	}

}

func reverseInts(input []ProfilMonster) []ProfilMonster {
	if len(input) == 0 {
		return input
	}
	return append(reverseInts(input[1:]), input[0])
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	fmt.Println(r.Body)
	if r.StatusCode== 404{
		return &errorString{"Could not found ressource"}
	}
	return json.NewDecoder(r.Body).Decode(target)
}

func GetSkill(s int) (Skill, error){
	var skill Skill
	url := "https://swarfarm.com/api/v2/skills/"
	errJson := getJson(url+strconv.Itoa(s), &skill)
	fmt.Println(url+strconv.Itoa(s))
	if errJson != nil {
		fmt.Println(errJson)
		return skill, &errorString{"skill not found"}
	}
	return skill, nil
}

func GetMonsterFile(){
	jsonFile, err := os.Open("monsters.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &MonstersList)
}

func GetRuneFile(){
	jsonFile, err := os.Open("runes.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &RunesMapping)
}

func GetMonsterId(mob string) int{
	for key, value := range MonstersList {
		if value[0] == mob {
			return key
		}
	}
	return 0
}

func GetMonster(mob string) (Monster, error) {
	var monsters Monster
	GetMonsterFile()
	////splittedString := strings.Split(m.Content, " ")
	url := "https://swarfarm.com/api/v2/monsters/"
	//splittedString := strings.TrimPrefix(m.Content, "?monster ")
	monsterId := GetMonsterId(mob)
	if monsterId > 0{
		errJson := getJson(url+strconv.Itoa(monsterId), &monsters)
		if errJson != nil {
			fmt.Println(errJson)
		}
		fmt.Println("monster found")
		return monsters, nil
	}else{
		return monsters,  &errorString{"monster not found"}
	}
}

func getBuildings(invoc string) (ProfilBuilding, error){
	var buildings ProfilBuilding
	url := "https://swarfarm.com/api/v2/profiles/"+invoc+"/buildings/"
	errJson := getJson(url, &buildings)
	if errJson != nil {
		fmt.Println(errJson)
		return buildings,  &errorString{"Summoners not found"}
	}
	return buildings, nil
}

func (e *errorString) Error() string {
	return e.s
}
func (n *Nstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

func (n *Skill) Print() (*discordgo.MessageEmbed) {
	//var result *discordgo.MessageEmbed
	fields := make(map[string]string)
	if n.Formula != "" {
		fields["Formula"] = n.Formula
	}
	for i:=0; i<len(n.SkillUp) ; i++  {
		fields["Skill Up"] = fields["Skill Up"] +"\nLv."+strconv.Itoa(i +1)+n.SkillUp[i]
	}
	result := CreateAnEmbed(n.Description, fields, n.Name, n.Image, "skills")
	//result = n.Name+" \n Description :"+n.Description +" \n Damage formula : " + n.Formula + " \n CoolTime "+ strconv.Itoa(n.CoolTime) +" \n Hits :"+strconv.Itoa(n.Hits)+" targets"
	return result
}

func (n *ProfilMonster) Print(name string) (string) {
	var result string

	result = name+" | "+strconv.Itoa(n.Star) +"* | lvl  " + strconv.Itoa(n.Level) +
		" \n Skill 1 Lvl : "+ strconv.Itoa(n.Skill1Level)+
		" \n Skill 2 Lvl : "+ strconv.Itoa(n.Skill2Level)+
		" \n Skill 3 Lvl : "+ strconv.Itoa(n.Skill3Level)+
		" \n Skill 4 Lvl : "+ strconv.Itoa(n.Skill4Level)+
		" \n Total HP : "+ strconv.Itoa(n.BaseHP+n.RuneHP)+
		" \n Total DEF : "+ strconv.Itoa(n.BaseDEF+n.RuneDEF)+
		" \n Total ACC : "+ strconv.Itoa(n.BaseACC+n.RuneACC)+
		" \n Total ATK : "+ strconv.Itoa(n.BaseATK+n.RuneATK)+
		" \n Total CR : "+ strconv.Itoa(n.BaseCR+n.RuneCR)+
		" \n Total CD : "+ strconv.Itoa(n.BaseCD+n.RuneCD)+
		" \n Total RES : "+ strconv.Itoa(n.BaseRES+n.RuneRES)+
		" \n Total SPD : "+ strconv.Itoa(n.BaseSPD+n.RuneSPD)+
		" \n Average Rune Efficiency : "+ fmt.Sprintf("%f", n.AverageRuneffeciency)+"%"

	return result
}

func (n *ProfilBuilding) FormatPrint() (map[string]string, string, string){
	fields := make(map[string]string)
	var name string
	//var image string
	var AllBuildings Buildings
	url := "https://swarfarm.com/api/v2/buildings/"
	errJson := getJson(url, &AllBuildings)
	if errJson != nil {
		fmt.Println(errJson)
	}
	for i:=0; i< len(n.Results);i++  {
		fields[AllBuildings.Results[n.Results[i].BuildingID -1].Name] = "**Level:** "+ strconv.Itoa(n.Results[i].Level)+"/"+strconv.Itoa(AllBuildings.Results[n.Results[i].BuildingID -1].MaxLevel)+
			"% \n **Bonus:** "+ AllBuildings.Results[n.Results[i].BuildingID -1].AffectedStat+ strconv.Itoa(AllBuildings.Results[n.Results[i].BuildingID -1].StatBonus[n.Results[i].Level -1])+ " in "+AllBuildings.Results[n.Results[i].BuildingID -1].Area
	}
	return fields, name, ""
}

func (n *ProfilMonster) FormatPrint() (map[string]string, string, string, string){
	GetRuneFile()
	fields := make(map[string]string)
	var name string
	var stars string
	var level string
	//var TempRune []Rune
	//TempRune = n.Runes
	RuneCheck := make(map[int]int)
	fields["Principales"] ="**HP:** "+strconv.Itoa(n.BaseHP+n.RuneHP)+
		"\n **ATK:** "+ strconv.Itoa(n.BaseATK+n.RuneATK)+
		"\n **DEF:** "+strconv.Itoa(n.BaseDEF+n.RuneDEF)+
		"\n **SPD:** "+strconv.Itoa(n.BaseSPD+n.RuneSPD)

	fields["Secondaires"] ="**CR:** "+strconv.Itoa(n.BaseCR+n.RuneCR)+
		"\n **CD:** "+strconv.Itoa(n.BaseCD+n.RuneCD)+
		"\n **RES:** "+strconv.Itoa(n.BaseRES+n.RuneRES)+
		"\n **ACC:** "+strconv.Itoa(n.BaseACC+n.RuneACC)
	fields["Efficiency"] =fmt.Sprintf("%f", n.AverageRuneffeciency)+"%"

	for i:=0; i< len(n.Runes) ; i++  {
		RuneCheck[n.Runes[i].Type] = RuneCheck[n.Runes[i].Type] + 1
	}
	IsBroken := false
	for key, value := range RuneCheck  {
		q, r := divmod(value, RunesMapping[key - 1].Set)
		for i:=0;i<q ;i++  {
			if fields["Set"] != "" {
				fields["Set"] = fields["Set"] + "/"
			}
			fields["Set"] = fields["Set"] + RunesMapping[key - 1].Type
		}
		if r>0 {
			IsBroken = true
		}
	}
	if IsBroken{
		if fields["Set"] != "" {
			fields["Set"] = fields["Set"] + "/"
		}
		fields["Set"] = fields["Set"] + "Broken"
	}

	if len(n.Runes)>1 {
		fields["Build"] = AttributeMapping[n.Runes[1].Main -1]
	}else{
		fields["Build"] = fields["Build"]+"None"
	}
	if len(n.Runes)>2 {
		fields["Build"] = fields["Build"]+"/"+AttributeMapping[n.Runes[3].Main -1]
	}else{
		fields["Build"] = fields["Build"]+"/None"
	}
	if len(n.Runes)>4 {
		fields["Build"] = fields["Build"]+"/"+AttributeMapping[n.Runes[5].Main -1]
	}else{
		fields["Build"] = fields["Build"]+"/None"
	}
	name = n.Name
	stars = strconv.Itoa(n.Star)
	level = strconv.Itoa(n.Level)
	return fields, name, stars, level
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	fmt.Println("num : ",numerator," deno : ",denominator)
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	fmt.Println("Quotient : ",quotient, " Remainder : ",remainder)
	return quotient, remainder
}

func (n *LeaderSkill) Print() (string) {
	var result string

	result = "Increase " + n.Attribute + " in "+ n.Area +" by "+strconv.Itoa(n.Amount)+"%"
	return result
}

func (n *Monster) PrintMonster() (string) {
	var result string
	fmt.Println(n.NaturalStars)

	result = n.Name+" | Element :"+n.Element +" | Type : " + n.Archetype + " | Nat "+ strconv.Itoa(n.NaturalStars) +"*"
	return result
}

func (n *Profile) PrintMonster() (string) {
	var result string
	result = MonstersList[n.Monster][0]+" "+ strconv.Itoa(n.Stars) +"*" + " lvl : "+ strconv.Itoa(n.Level)
	return result
}

func CreateAnEmbed(description string, fields map[string]string, title string, img string, which string) *discordgo.MessageEmbed{
	var fieldsToSend []*discordgo.MessageEmbedField
	for key, value := range fields {
		fieldsToSend= append(fieldsToSend,&discordgo.MessageEmbedField{
					Name:   key,
					Value:  value,
					Inline: false,
				} )
	}
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: description,
		Fields: fieldsToSend,
		//Fields: []*discordgo.MessageEmbedField{
		//	&discordgo.MessageEmbedField{
		//		Name:   "I am a field",
		//		Value:  "I am a value",
		//		Inline: true,
		//	},
		//	&discordgo.MessageEmbedField{
		//		Name:   "I am a second field",
		//		Value:  "I am a value",
		//		Inline: true,
		//	},
		//},
		//Image: &discordgo.MessageEmbedImage{
		//	URL: "https://swarfarm.com/static/herders/images/monsters/"+img,
		//},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://swarfarm.com/static/herders/images/"+which+"/"+img,
		},
		//Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     title,
	}
	fmt.Println(embed.Thumbnail)
	return embed
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}