package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/xaviatwork/ipam/ipamautopilot"
)

type GpsIpam struct {
	Source string
}

func (gpsipam GpsIpam) RangeById(id int) (*ipamautopilot.Range, error) {
	lzrange := &ipamautopilot.Range{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/ranges/%d", gpsipam.Source, id))
	if err != nil {
		return lzrange, err
	}
	if err := json.Unmarshal(b, &lzrange); err != nil {
		return lzrange, err
	}
	return lzrange, nil
}
func (gpsipam GpsIpam) Ranges() (*[]ipamautopilot.Range, error) {
	ranges := &[]ipamautopilot.Range{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/ranges", gpsipam.Source))
	if err != nil {
		return ranges, err
	}
	if err := json.Unmarshal(b, &ranges); err != nil {
		return ranges, err
	}
	return ranges, nil
}
func (gpsipam GpsIpam) RoutingDomainById(id int) (*ipamautopilot.RoutingDomain, error) {
	routingdomain := &ipamautopilot.RoutingDomain{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/domains/%d", gpsipam.Source, id))
	if err != nil {
		return routingdomain, err
	}
	if err := json.Unmarshal(b, &routingdomain); err != nil {
		return routingdomain, err
	}
	return routingdomain, nil
}
func (gpsipam GpsIpam) RoutingDomains() (*[]ipamautopilot.RoutingDomain, error) {
	domains := &[]ipamautopilot.RoutingDomain{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/domains", gpsipam.Source))
	if err != nil {
		return domains, err
	}
	if err := json.Unmarshal(b, &domains); err != nil {
		return domains, err
	}
	return domains, nil
}

func (gpsipam GpsIpam) getToken() string {
	return os.Getenv("IPAM_TOKEN")
}

func (gpsipam GpsIpam) doRequest(url string) ([]byte, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("Authorization", "bearer "+gpsipam.getToken())

	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return []byte{}, fmt.Errorf("http error %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (gpsipam GpsIpam) Status() error {
	b, err := gpsipam.doRequest(gpsipam.Source)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func GetRangeById(ipam ipamautopilot.Ipam, opts Opts) {
	iprange, err := ipam.RangeById(opts.Id)
	iprng := Pokemonize(iprange)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	if opts.Pretty {
		fmt.Printf("%s\n", iprng.PrettyString())
		return
	}
	fmt.Printf("%s\n", iprng.String())
}

func GetRangesWithParent(ipam ipamautopilot.Ipam, opts Opts) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		pkr := Pokemonize(&r)
		if pkr.Parent_id == opts.Parent {
			if opts.Pretty {
				fmt.Printf("%s\n", pkr.PrettyString())
				continue
			}
			fmt.Printf("%s", pkr.String())
		}
	}
}

func SearchStringInRanges(ipam ipamautopilot.Ipam, opts Opts) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		pkr := Pokemonize(&r)
		if searchString(opts.SearchString, pkr.Name, pkr.Cidr) {
			if opts.Pretty {
				fmt.Printf("%s\n", pkr.PrettyString())
				continue
			}
			fmt.Printf("%s", pkr.String())
		}
	}
}

func GetDomainById(ipam ipamautopilot.Ipam, opts Opts) {
	domain, err := ipam.RoutingDomainById(opts.Id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	if opts.Pretty {
		fmt.Printf("%s\n", domain.PrettyString())
		return
	}
	fmt.Printf("%s\n", domain.String())
}

func GetNonAllocatedIPs(ipam ipamautopilot.Ipam, opts Opts) {
	iprange, err := ipam.RangeById(opts.Free)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	pokeRange := Pokemonize(iprange)

	if pokeRange.Parent_id != -1 {
		fmt.Printf("Range %d (name: %s) is not a main range\n", opts.Free, pokeRange.Name)
		os.Exit(1)
	}

	// lleig
	s := strings.SplitAfter(pokeRange.Cidr, "/")
	strblock := s[len(s)-1]
	size, _ := strconv.Atoi(strblock)
	ipsOnRange := math.Pow(2, float64(32)-float64(size))

	subranges, err := ipam.Ranges()
	if err != nil {
		log.Printf("error: %s", err.Error())
	}
	available := ipsOnRange
	var allocated float64

	pokerange := Pokemonize(&pokeRange)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"CIDR", "Allocated IPs", "Available IPs", "CIDR Name"})

	t.AppendRow(table.Row{pokerange.Cidr, 0, available, pokerange.Name})
	t.AppendSeparator()

	// fmt.Printf("total IPs: %.0f, CIDR: %s (%q))\n", ipsOnRange, iprange.Cidr, iprange.Name)
	for _, sr := range *subranges {
		if sr.Parent_id == opts.Free {
			pokeSubnet := Pokemonize(&sr)
			s := strings.SplitAfter(pokeSubnet.Cidr, "/")
			strblock := s[len(s)-1]
			size, _ := strconv.Atoi(strblock)
			ipsOnSubnet := math.Pow(2, float64(32)-float64(size))
			// fmt.Printf("ips on subnet %.0f (CIDR: %s)\n", ipsOnSubnet, sr.Cidr)
			available = available - ipsOnSubnet
			allocated = allocated + ipsOnSubnet
			// fmt.Printf("Available %4.0f (after allocating %3.0f IPs, %s)\n", available, ipsOnSubnet, sr.Cidr)

			t.AppendRow(table.Row{pokeSubnet.Cidr, ipsOnSubnet, available, pokeSubnet.Name})
		}
	}
	t.AppendSeparator()
	t.AppendRow(table.Row{"Total", allocated, available, ""})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func SearchStringInDomains(ipam ipamautopilot.Ipam, opts Opts) {
	domains, err := ipam.RoutingDomains()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, d := range *domains {
		if searchString(opts.SearchString, d.Name, d.Vpcs) {
			if opts.Pretty {
				fmt.Printf("%s\n", d.PrettyString())
				continue
			}
			fmt.Printf("%s", d.String())
		}
	}
}

// SearchString returns true if any of the ss[1:] strings contains ss[0]
func searchString(ss ...string) bool {
	searchString := strings.ToLower(ss[0])
	found := false
	for _, s := range ss[1:] {
		if strings.Contains(strings.ToLower(s), searchString) {
			found = true
		}
	}
	return found
}

func Pokemonize(r *ipamautopilot.Range) ipamautopilot.Range {
	if os.Getenv("POKEMONIZE") == "yes" {
		pokemons := []string{"Bulbasaur", "Ivysaur", "Venusaur", "Charmander", "Charmeleon", "Charizard", "Squirtle", "Wartortle", "Blastoise", "Caterpie", "Metapod", "Butterfree", "Weedle", "Kakuna", "Beedrill", "Pidgey", "Pidgeotto", "Pidgeot", "Rattata", "Raticate", "Spearow", "Fearow", "Ekans", "Arbok", "Pikachu", "Raichu", "Sandshrew", "Sandslash", "Nidoran♀", "Nidorina", "Nidoqueen", "Nidoran♂", "Nidorino", "Nidoking", "Clefairy", "Clefable", "Vulpix", "Ninetales", "Jigglypuff", "Wigglytuff", "Zubat", "Golbat", "Oddish", "Gloom", "Vileplume", "Paras", "Parasect", "Venonat", "Venomoth", "Diglett", "Dugtrio", "Meowth", "Persian", "Psyduck", "Golduck", "Mankey", "Primeape", "Growlithe", "Arcanine", "Poliwag", "Poliwhirl", "Poliwrath", "Abra", "Kadabra", "Alakazam", "Machop", "Machoke", "Machamp", "Bellsprout", "Weepinbell", "Victreebel", "Tentacool", "Tentacruel", "Geodude", "Graveler", "Golem", "Ponyta", "Rapidash", "Slowpoke", "Slowbro", "Magnemite", "Magneton", "Farfetch’d", "Doduo", "Dodrio", "Seel", "Dewgong", "Grimer", "Muk", "Shellder", "Cloyster", "Gastly", "Haunter", "Gengar", "Onix", "Drowzee", "Hypno", "Krabby", "Kingler", "Voltorb", "Electrode", "Exeggcute", "Exeggutor", "Cubone", "Marowak", "Hitmonlee", "Hitmonchan", "Lickitung", "Koffing", "Weezing", "Rhyhorn", "Rhydon", "Chansey", "Tangela", "Kangaskhan", "Horsea", "Seadra", "Goldeen", "Seaking", "Staryu", "Starmie", "Mr._Mime", "Scyther", "Jynx", "Electabuzz", "Magmar", "Pinsir", "Tauros", "Magikarp", "Gyarados", "Lapras", "Ditto", "Eevee", "Vaporeon", "Jolteon", "Flareon", "Porygon", "Omanyte", "Omastar", "Kabuto", "Kabutops", "Aerodactyl", "Snorlax", "Articuno", "Zapdos", "Moltres", "Dratini", "Dragonair", "Dragonite", "Mewtwo", "Mew", "Chikorita", "Bayleef", "Meganium", "Cyndaquil", "Quilava", "Typhlosion", "Totodile", "Croconaw", "Feraligatr", "Sentret", "Furret", "Hoothoot", "Noctowl", "Ledyba", "Ledian", "Spinarak", "Ariados", "Crobat", "Chinchou", "Lanturn", "Pichu", "Cleffa", "Igglybuff", "Togepi", "Togetic", "Natu", "Xatu", "Mareep", "Flaaffy", "Ampharos", "Bellossom", "Marill", "Azumarill", "Sudowoodo", "Politoed", "Hoppip", "Skiploom", "Jumpluff", "Aipom", "Sunkern", "Sunflora", "Yanma", "Wooper", "Quagsire", "Espeon", "Umbreon", "Murkrow", "Slowking", "Misdreavus", "Unown", "Wobbuffet", "Girafarig", "Pineco", "Forretress", "Dunsparce", "Gligar", "Steelix", "Snubbull", "Granbull", "Qwilfish", "Scizor", "Shuckle", "Heracross", "Sneasel", "Teddiursa", "Ursaring", "Slugma", "Magcargo", "Swinub", "Piloswine", "Corsola", "Remoraid", "Octillery", "Delibird", "Mantine", "Skarmory", "Houndour", "Houndoom", "Kingdra", "Phanpy", "Donphan", "Porygon2", "Stantler", "Smeargle", "Tyrogue", "Hitmontop", "Smoochum", "Elekid", "Magby", "Miltank", "Blissey", "Raikou", "Entei", "Suicune", "Larvitar", "Pupitar", "Tyranitar", "Lugia", "Ho-Oh", "Celebi", "Treecko", "Grovyle", "Sceptile", "Torchic", "Combusken", "Blaziken", "Mudkip", "Marshtomp", "Swampert", "Poochyena", "Mightyena", "Zigzagoon", "Linoone", "Wurmple", "Silcoon", "Beautifly", "Cascoon", "Dustox", "Lotad", "Lombre", "Ludicolo", "Seedot", "Nuzleaf", "Shiftry", "Taillow", "Swellow", "Wingull", "Pelipper", "Ralts", "Kirlia", "Gardevoir", "Surskit", "Masquerain", "Shroomish", "Breloom", "Slakoth", "Vigoroth", "Slaking", "Nincada", "Ninjask", "Shedinja", "Whismur", "Loudred", "Exploud", "Makuhita", "Hariyama", "Azurill", "Nosepass", "Skitty", "Delcatty", "Sableye", "Mawile", "Aron", "Lairon", "Aggron", "Meditite", "Medicham", "Electrike", "Manectric", "Plusle", "Minun", "Volbeat", "Illumise", "Roselia", "Gulpin", "Swalot", "Carvanha", "Sharpedo", "Wailmer", "Wailord", "Numel", "Camerupt", "Torkoal", "Spoink", "Grumpig", "Spinda", "Trapinch", "Vibrava", "Flygon", "Cacnea", "Cacturne", "Swablu", "Altaria", "Zangoose", "Seviper", "Lunatone", "Solrock", "Barboach", "Whiscash", "Corphish", "Crawdaunt", "Baltoy", "Claydol", "Lileep", "Cradily", "Anorith", "Armaldo", "Feebas", "Milotic", "Castform", "Kecleon", "Shuppet", "Banette", "Duskull", "Dusclops", "Tropius", "Chimecho", "Absol", "Wynaut", "Snorunt", "Glalie", "Spheal", "Sealeo", "Walrein", "Clamperl", "Huntail", "Gorebyss", "Relicanth", "Luvdisc", "Bagon", "Shelgon", "Salamence", "Beldum", "Metang", "Metagross", "Regirock", "Regice", "Registeel", "Latias", "Latios", "Kyogre", "Groudon", "Rayquaza", "Jirachi", "Deoxys", "Turtwig", "Grotle", "Torterra", "Chimchar", "Monferno", "Infernape", "Piplup", "Prinplup", "Empoleon", "Starly", "Staravia", "Staraptor", "Bidoof", "Bibarel", "Kricketot", "Kricketune", "Shinx", "Luxio", "Luxray", "Budew", "Roserade", "Cranidos", "Rampardos", "Shieldon", "Bastiodon", "Burmy", "Wormadam", "Mothim", "Combee", "Vespiquen", "Pachirisu", "Buizel", "Floatzel", "Cherubi", "Cherrim", "Shellos", "Gastrodon", "Ambipom", "Drifloon", "Drifblim", "Buneary", "Lopunny", "Mismagius", "Honchkrow", "Glameow", "Purugly", "Chingling", "Stunky", "Skuntank", "Bronzor", "Bronzong", "Bonsly", "Mime_Jr.", "Happiny", "Chatot", "Spiritomb", "Gible", "Gabite", "Garchomp", "Munchlax", "Riolu", "Lucario", "Hippopotas", "Hippowdon", "Skorupi", "Drapion", "Croagunk", "Toxicroak", "Carnivine", "Finneon", "Lumineon", "Mantyke", "Snover", "Abomasnow", "Weavile", "Magnezone", "Lickilicky", "Rhyperior", "Tangrowth", "Electivire", "Magmortar", "Togekiss", "Yanmega", "Leafeon", "Glaceon", "Gliscor", "Mamoswine", "Porygon-Z", "Gallade", "Probopass", "Dusknoir", "Froslass", "Rotom", "Uxie", "Mesprit", "Azelf", "Dialga", "Palkia", "Heatran", "Regigigas", "Giratina", "Cresselia", "Phione", "Manaphy", "Darkrai", "Shaymin", "Arceus", "Victini", "Snivy", "Servine", "Serperior", "Tepig", "Pignite", "Emboar", "Oshawott", "Dewott", "Samurott", "Patrat", "Watchog", "Lillipup", "Herdier", "Stoutland", "Purrloin", "Liepard", "Pansage", "Simisage", "Pansear", "Simisear", "Panpour", "Simipour", "Munna", "Musharna", "Pidove", "Tranquill", "Unfezant", "Blitzle", "Zebstrika", "Roggenrola", "Boldore", "Gigalith", "Woobat", "Swoobat", "Drilbur", "Excadrill", "Audino", "Timburr", "Gurdurr", "Conkeldurr", "Tympole", "Palpitoad", "Seismitoad", "Throh", "Sawk", "Sewaddle", "Swadloon", "Leavanny", "Venipede", "Whirlipede", "Scolipede", "Cottonee", "Whimsicott", "Petilil", "Lilligant", "Basculin", "Sandile", "Krokorok", "Krookodile", "Darumaka", "Darmanitan", "Maractus", "Dwebble", "Crustle", "Scraggy", "Scrafty", "Sigilyph", "Yamask", "Cofagrigus", "Tirtouga", "Carracosta", "Archen", "Archeops", "Trubbish", "Garbodor", "Zorua", "Zoroark", "Minccino", "Cinccino", "Gothita", "Gothorita", "Gothitelle", "Solosis", "Duosion", "Reuniclus", "Ducklett", "Swanna", "Vanillite", "Vanillish", "Vanilluxe", "Deerling", "Sawsbuck", "Emolga", "Karrablast", "Escavalier", "Foongus", "Amoonguss", "Frillish", "Jellicent", "Alomomola", "Joltik", "Galvantula", "Ferroseed", "Ferrothorn", "Klink", "Klang", "Klinklang", "Tynamo", "Eelektrik", "Eelektross", "Elgyem", "Beheeyem", "Litwick", "Lampent", "Chandelure", "Axew", "Fraxure", "Haxorus", "Cubchoo", "Beartic", "Cryogonal", "Shelmet", "Accelgor", "Stunfisk", "Mienfoo", "Mienshao", "Druddigon", "Golett", "Golurk", "Pawniard", "Bisharp", "Bouffalant", "Rufflet", "Braviary", "Vullaby", "Mandibuzz", "Heatmor", "Durant", "Deino", "Zweilous", "Hydreigon", "Larvesta", "Volcarona", "Cobalion", "Terrakion", "Virizion", "Tornadus", "Thundurus", "Reshiram", "Zekrom", "Landorus", "Kyurem", "Keldeo", "Meloetta", "Genesect", "Chespin", "Quilladin", "Chesnaught", "Fennekin", "Braixen", "Delphox", "Froakie", "Frogadier", "Greninja", "Bunnelby", "Diggersby", "Fletchling", "Fletchinder", "Talonflame", "Scatterbug", "Spewpa", "Vivillon", "Litleo", "Pyroar", "Flabébé", "Floette", "Florges", "Skiddo", "Gogoat", "Pancham", "Pangoro", "Furfrou", "Espurr", "Meowstic", "Honedge", "Doublade", "Aegislash", "Spritzee", "Aromatisse", "Swirlix", "Slurpuff", "Inkay", "Malamar", "Binacle", "Barbaracle", "Skrelp", "Dragalge", "Clauncher", "Clawitzer", "Helioptile", "Heliolisk", "Tyrunt", "Tyrantrum", "Amaura", "Aurorus", "Sylveon", "Hawlucha", "Dedenne", "Carbink", "Goomy", "Sliggoo", "Goodra", "Klefki", "Phantump", "Trevenant", "Pumpkaboo", "Gourgeist", "Bergmite", "Avalugg", "Noibat", "Noivern", "Xerneas", "Yveltal", "Zygarde", "Diancie", "Hoopa", "Volcanion", "Rowlet", "Dartrix", "Decidueye", "Litten", "Torracat", "Incineroar", "Popplio", "Brionne", "Primarina", "Pikipek", "Trumbeak", "Toucannon", "Yungoos", "Gumshoos", "Grubbin", "Charjabug", "Vikavolt", "Crabrawler", "Crabominable", "Oricorio", "Cutiefly", "Ribombee", "Rockruff", "Lycanroc", "Wishiwashi", "Mareanie", "Toxapex", "Mudbray", "Mudsdale", "Dewpider", "Araquanid", "Fomantis", "Lurantis", "Morelull", "Shiinotic", "Salandit", "Salazzle", "Stufful", "Bewear", "Bounsweet", "Steenee", "Tsareena", "Comfey", "Oranguru", "Passimian", "Wimpod", "Golisopod", "Sandygast", "Palossand", "Pyukumuku", "Type:_Null", "Silvally", "Minior", "Komala", "Turtonator", "Togedemaru", "Mimikyu", "Bruxish", "Drampa", "Dhelmise", "Jangmo-o", "Hakamo-o", "Kommo-o", "Tapu_Koko", "Tapu_Lele", "Tapu_Bulu", "Tapu_Fini", "Cosmog", "Cosmoem", "Solgaleo", "Lunala", "Nihilego", "Buzzwole", "Pheromosa", "Xurkitree", "Celesteela", "Kartana", "Guzzlord", "Necrozma", "Magearna", "Marshadow", "Poipole", "Naganadel", "Stakataka", "Blacephalon", "Zeraora", "Meltan", "Melmetal"}
		s := strings.SplitAfter(r.Cidr, "/")
		cidrBlock := s[len(s)-1]
		x, y := rand.IntN(255), rand.IntN(255)
		r.Cidr = fmt.Sprintf("192.168.%d.%d/%s", x, y, cidrBlock)
		r.Name = fmt.Sprintf("%s-%s-%s", pokemons[rand.IntN(len(pokemons))], pokemons[rand.IntN(len(pokemons))], pokemons[rand.IntN(len(pokemons))])
		return *r
	}
	return *r
}
