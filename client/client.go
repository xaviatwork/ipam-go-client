package client

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xaviatwork/ipam/ipamautopilot"
)

func GetRangeById(ipam ipamautopilot.Ipam, opts Opts) {
	iprange, err := ipam.RangeById(opts.Id)
	iprng := Anonymize(iprange)
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
		pkr := Anonymize(&r)
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
		pkr := Anonymize(&r)
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

func ipsOnRange(cidr string) float64 {
	_, mask, _ := net.ParseCIDR(cidr)
	size, _ := mask.Mask.Size()
	return math.Pow(2, float64(32)-float64(size))
}

func GetNonAllocatedIPs(ipam ipamautopilot.Ipam, opts Opts) {
	mr, err := ipam.RangeById(opts.Id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	mainRange := Anonymize(mr)

	if mainRange.Parent_id != -1 {
		fmt.Printf("Range %d (name: %s) is not a main range\n", opts.Id, mainRange.Name)
		os.Exit(1)
	}

	availableIPs := ipsOnRange(mainRange.Cidr)
	sizeMainRange := availableIPs

	srs, err := ipam.Ranges()
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	var allocated float64
	switch opts.Format {
	case "table":
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"CIDR", "Allocated IPs", "Available IPs", "CIDR Name"})
		t.AppendRow(table.Row{mainRange.Cidr, 0, availableIPs, mainRange.Name})
		t.AppendSeparator()

		for _, sr := range *srs {
			if sr.Parent_id == opts.Id {
				subnetRange := Anonymize(&sr)
				ipsSubnet := ipsOnRange(subnetRange.Cidr)

				availableIPs = availableIPs - ipsSubnet
				allocated = allocated + ipsSubnet
				t.AppendRow(table.Row{subnetRange.Cidr, ipsSubnet, availableIPs, subnetRange.Name})
			}
		}
		t.AppendSeparator()
		t.AppendRow(table.Row{"Total", allocated, availableIPs, ""})
		t.SetStyle(table.StyleLight)
		t.Render()
	case "number":
		for _, sr := range *srs {
			if sr.Parent_id == opts.Id {
				subnetRange := Anonymize(&sr)
				ipsSubnet := ipsOnRange(subnetRange.Cidr)

				availableIPs = availableIPs - ipsSubnet
				allocated = allocated + ipsSubnet
			}
		}
		fmt.Printf("%.0f", availableIPs)
	case "json":
		for _, sr := range *srs {
			if sr.Parent_id == opts.Id {
				subnetRange := Anonymize(&sr)
				ipsSubnet := ipsOnRange(subnetRange.Cidr)

				availableIPs = availableIPs - ipsSubnet
				allocated = allocated + ipsSubnet
			}
		}
		fmt.Printf(`{"name": "%s", "range": "%s", "id": %d, "ip addreses": {"total": %.0f, "allocated": %.0f, "available": %.0f}}`, mainRange.Name, mainRange.Cidr, mainRange.Subnet_id, sizeMainRange, allocated, availableIPs)
	}

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

func Anonymize(r *ipamautopilot.Range) ipamautopilot.Range {
	if os.Getenv("POKEMONIZE") == "true" || os.Getenv("ANONYMIZE") == "true" {
		pokemons := []string{"Bulbasaur", "Ivysaur", "Venusaur", "Charmander", "Charmeleon", "Charizard", "Squirtle", "Wartortle", "Blastoise", "Caterpie", "Metapod", "Butterfree", "Weedle", "Kakuna", "Beedrill", "Pidgey", "Pidgeotto", "Pidgeot", "Rattata", "Raticate", "Spearow", "Fearow", "Ekans", "Arbok", "Pikachu", "Raichu", "Sandshrew", "Sandslash", "Nidoran♀", "Nidorina", "Nidoqueen", "Nidoran♂", "Nidorino", "Nidoking", "Clefairy", "Clefable", "Vulpix", "Ninetales", "Jigglypuff", "Wigglytuff", "Zubat", "Golbat", "Oddish", "Gloom", "Vileplume", "Paras", "Parasect", "Venonat", "Venomoth", "Diglett", "Dugtrio", "Meowth", "Persian", "Psyduck", "Golduck", "Mankey", "Primeape", "Growlithe", "Arcanine", "Poliwag", "Poliwhirl", "Poliwrath", "Abra", "Kadabra", "Alakazam", "Machop", "Machoke", "Machamp", "Bellsprout", "Weepinbell", "Victreebel", "Tentacool", "Tentacruel", "Geodude", "Graveler", "Golem", "Ponyta", "Rapidash", "Slowpoke", "Slowbro", "Magnemite", "Magneton", "Farfetch’d", "Doduo", "Dodrio", "Seel", "Dewgong", "Grimer", "Muk", "Shellder", "Cloyster", "Gastly", "Haunter", "Gengar", "Onix", "Drowzee", "Hypno", "Krabby", "Kingler", "Voltorb", "Electrode", "Exeggcute", "Exeggutor", "Cubone", "Marowak", "Hitmonlee", "Hitmonchan", "Lickitung", "Koffing", "Weezing", "Rhyhorn", "Rhydon", "Chansey", "Tangela", "Kangaskhan", "Horsea", "Seadra", "Goldeen", "Seaking", "Staryu", "Starmie", "Mr._Mime", "Scyther", "Jynx", "Electabuzz", "Magmar", "Pinsir", "Tauros", "Magikarp", "Gyarados", "Lapras", "Ditto", "Eevee", "Vaporeon", "Jolteon", "Flareon", "Porygon", "Omanyte", "Omastar", "Kabuto", "Kabutops", "Aerodactyl", "Snorlax", "Articuno", "Zapdos", "Moltres", "Dratini", "Dragonair", "Dragonite", "Mewtwo", "Mew", "Chikorita", "Bayleef", "Meganium", "Cyndaquil", "Quilava", "Typhlosion", "Totodile", "Croconaw", "Feraligatr", "Sentret", "Furret", "Hoothoot", "Noctowl", "Ledyba", "Ledian", "Spinarak", "Ariados", "Crobat", "Chinchou", "Lanturn", "Pichu", "Cleffa", "Igglybuff", "Togepi", "Togetic", "Natu", "Xatu", "Mareep", "Flaaffy", "Ampharos", "Bellossom", "Marill", "Azumarill", "Sudowoodo", "Politoed", "Hoppip", "Skiploom", "Jumpluff", "Aipom", "Sunkern", "Sunflora", "Yanma", "Wooper", "Quagsire", "Espeon", "Umbreon", "Murkrow", "Slowking", "Misdreavus", "Unown", "Wobbuffet", "Girafarig", "Pineco", "Forretress", "Dunsparce", "Gligar", "Steelix", "Snubbull", "Granbull", "Qwilfish", "Scizor", "Shuckle", "Heracross", "Sneasel", "Teddiursa", "Ursaring", "Slugma", "Magcargo", "Swinub", "Piloswine", "Corsola", "Remoraid", "Octillery", "Delibird", "Mantine", "Skarmory", "Houndour", "Houndoom", "Kingdra", "Phanpy", "Donphan", "Porygon2", "Stantler", "Smeargle", "Tyrogue", "Hitmontop", "Smoochum", "Elekid", "Magby", "Miltank", "Blissey", "Raikou", "Entei", "Suicune", "Larvitar", "Pupitar", "Tyranitar", "Lugia", "Ho-Oh", "Celebi", "Treecko", "Grovyle", "Sceptile", "Torchic", "Combusken", "Blaziken", "Mudkip", "Marshtomp", "Swampert", "Poochyena", "Mightyena", "Zigzagoon", "Linoone", "Wurmple", "Silcoon", "Beautifly", "Cascoon", "Dustox", "Lotad", "Lombre", "Ludicolo", "Seedot", "Nuzleaf", "Shiftry", "Taillow", "Swellow", "Wingull", "Pelipper", "Ralts", "Kirlia", "Gardevoir", "Surskit", "Masquerain", "Shroomish", "Breloom", "Slakoth", "Vigoroth", "Slaking", "Nincada", "Ninjask", "Shedinja", "Whismur", "Loudred", "Exploud", "Makuhita", "Hariyama", "Azurill", "Nosepass", "Skitty", "Delcatty", "Sableye", "Mawile", "Aron", "Lairon", "Aggron", "Meditite", "Medicham", "Electrike", "Manectric", "Plusle", "Minun", "Volbeat", "Illumise", "Roselia", "Gulpin", "Swalot", "Carvanha", "Sharpedo", "Wailmer", "Wailord", "Numel", "Camerupt", "Torkoal", "Spoink", "Grumpig", "Spinda", "Trapinch", "Vibrava", "Flygon", "Cacnea", "Cacturne", "Swablu", "Altaria", "Zangoose", "Seviper", "Lunatone", "Solrock", "Barboach", "Whiscash", "Corphish", "Crawdaunt", "Baltoy", "Claydol", "Lileep", "Cradily", "Anorith", "Armaldo", "Feebas", "Milotic", "Castform", "Kecleon", "Shuppet", "Banette", "Duskull", "Dusclops", "Tropius", "Chimecho", "Absol", "Wynaut", "Snorunt", "Glalie", "Spheal", "Sealeo", "Walrein", "Clamperl", "Huntail", "Gorebyss", "Relicanth", "Luvdisc", "Bagon", "Shelgon", "Salamence", "Beldum", "Metang", "Metagross", "Regirock", "Regice", "Registeel", "Latias", "Latios", "Kyogre", "Groudon", "Rayquaza", "Jirachi", "Deoxys", "Turtwig", "Grotle", "Torterra", "Chimchar", "Monferno", "Infernape", "Piplup", "Prinplup", "Empoleon", "Starly", "Staravia", "Staraptor", "Bidoof", "Bibarel", "Kricketot", "Kricketune", "Shinx", "Luxio", "Luxray", "Budew", "Roserade", "Cranidos", "Rampardos", "Shieldon", "Bastiodon", "Burmy", "Wormadam", "Mothim", "Combee", "Vespiquen", "Pachirisu", "Buizel", "Floatzel", "Cherubi", "Cherrim", "Shellos", "Gastrodon", "Ambipom", "Drifloon", "Drifblim", "Buneary", "Lopunny", "Mismagius", "Honchkrow", "Glameow", "Purugly", "Chingling", "Stunky", "Skuntank", "Bronzor", "Bronzong", "Bonsly", "Mime_Jr.", "Happiny", "Chatot", "Spiritomb", "Gible", "Gabite", "Garchomp", "Munchlax", "Riolu", "Lucario", "Hippopotas", "Hippowdon", "Skorupi", "Drapion", "Croagunk", "Toxicroak", "Carnivine", "Finneon", "Lumineon", "Mantyke", "Snover", "Abomasnow", "Weavile", "Magnezone", "Lickilicky", "Rhyperior", "Tangrowth", "Electivire", "Magmortar", "Togekiss", "Yanmega", "Leafeon", "Glaceon", "Gliscor", "Mamoswine", "Porygon-Z", "Gallade", "Probopass", "Dusknoir", "Froslass", "Rotom", "Uxie", "Mesprit", "Azelf", "Dialga", "Palkia", "Heatran", "Regigigas", "Giratina", "Cresselia", "Phione", "Manaphy", "Darkrai", "Shaymin", "Arceus", "Victini", "Snivy", "Servine", "Serperior", "Tepig", "Pignite", "Emboar", "Oshawott", "Dewott", "Samurott", "Patrat", "Watchog", "Lillipup", "Herdier", "Stoutland", "Purrloin", "Liepard", "Pansage", "Simisage", "Pansear", "Simisear", "Panpour", "Simipour", "Munna", "Musharna", "Pidove", "Tranquill", "Unfezant", "Blitzle", "Zebstrika", "Roggenrola", "Boldore", "Gigalith", "Woobat", "Swoobat", "Drilbur", "Excadrill", "Audino", "Timburr", "Gurdurr", "Conkeldurr", "Tympole", "Palpitoad", "Seismitoad", "Throh", "Sawk", "Sewaddle", "Swadloon", "Leavanny", "Venipede", "Whirlipede", "Scolipede", "Cottonee", "Whimsicott", "Petilil", "Lilligant", "Basculin", "Sandile", "Krokorok", "Krookodile", "Darumaka", "Darmanitan", "Maractus", "Dwebble", "Crustle", "Scraggy", "Scrafty", "Sigilyph", "Yamask", "Cofagrigus", "Tirtouga", "Carracosta", "Archen", "Archeops", "Trubbish", "Garbodor", "Zorua", "Zoroark", "Minccino", "Cinccino", "Gothita", "Gothorita", "Gothitelle", "Solosis", "Duosion", "Reuniclus", "Ducklett", "Swanna", "Vanillite", "Vanillish", "Vanilluxe", "Deerling", "Sawsbuck", "Emolga", "Karrablast", "Escavalier", "Foongus", "Amoonguss", "Frillish", "Jellicent", "Alomomola", "Joltik", "Galvantula", "Ferroseed", "Ferrothorn", "Klink", "Klang", "Klinklang", "Tynamo", "Eelektrik", "Eelektross", "Elgyem", "Beheeyem", "Litwick", "Lampent", "Chandelure", "Axew", "Fraxure", "Haxorus", "Cubchoo", "Beartic", "Cryogonal", "Shelmet", "Accelgor", "Stunfisk", "Mienfoo", "Mienshao", "Druddigon", "Golett", "Golurk", "Pawniard", "Bisharp", "Bouffalant", "Rufflet", "Braviary", "Vullaby", "Mandibuzz", "Heatmor", "Durant", "Deino", "Zweilous", "Hydreigon", "Larvesta", "Volcarona", "Cobalion", "Terrakion", "Virizion", "Tornadus", "Thundurus", "Reshiram", "Zekrom", "Landorus", "Kyurem", "Keldeo", "Meloetta", "Genesect", "Chespin", "Quilladin", "Chesnaught", "Fennekin", "Braixen", "Delphox", "Froakie", "Frogadier", "Greninja", "Bunnelby", "Diggersby", "Fletchling", "Fletchinder", "Talonflame", "Scatterbug", "Spewpa", "Vivillon", "Litleo", "Pyroar", "Flabébé", "Floette", "Florges", "Skiddo", "Gogoat", "Pancham", "Pangoro", "Furfrou", "Espurr", "Meowstic", "Honedge", "Doublade", "Aegislash", "Spritzee", "Aromatisse", "Swirlix", "Slurpuff", "Inkay", "Malamar", "Binacle", "Barbaracle", "Skrelp", "Dragalge", "Clauncher", "Clawitzer", "Helioptile", "Heliolisk", "Tyrunt", "Tyrantrum", "Amaura", "Aurorus", "Sylveon", "Hawlucha", "Dedenne", "Carbink", "Goomy", "Sliggoo", "Goodra", "Klefki", "Phantump", "Trevenant", "Pumpkaboo", "Gourgeist", "Bergmite", "Avalugg", "Noibat", "Noivern", "Xerneas", "Yveltal", "Zygarde", "Diancie", "Hoopa", "Volcanion", "Rowlet", "Dartrix", "Decidueye", "Litten", "Torracat", "Incineroar", "Popplio", "Brionne", "Primarina", "Pikipek", "Trumbeak", "Toucannon", "Yungoos", "Gumshoos", "Grubbin", "Charjabug", "Vikavolt", "Crabrawler", "Crabominable", "Oricorio", "Cutiefly", "Ribombee", "Rockruff", "Lycanroc", "Wishiwashi", "Mareanie", "Toxapex", "Mudbray", "Mudsdale", "Dewpider", "Araquanid", "Fomantis", "Lurantis", "Morelull", "Shiinotic", "Salandit", "Salazzle", "Stufful", "Bewear", "Bounsweet", "Steenee", "Tsareena", "Comfey", "Oranguru", "Passimian", "Wimpod", "Golisopod", "Sandygast", "Palossand", "Pyukumuku", "Type:_Null", "Silvally", "Minior", "Komala", "Turtonator", "Togedemaru", "Mimikyu", "Bruxish", "Drampa", "Dhelmise", "Jangmo-o", "Hakamo-o", "Kommo-o", "Tapu_Koko", "Tapu_Lele", "Tapu_Bulu", "Tapu_Fini", "Cosmog", "Cosmoem", "Solgaleo", "Lunala", "Nihilego", "Buzzwole", "Pheromosa", "Xurkitree", "Celesteela", "Kartana", "Guzzlord", "Necrozma", "Magearna", "Marshadow", "Poipole", "Naganadel", "Stakataka", "Blacephalon", "Zeraora", "Meltan", "Melmetal"}
		_, mask, _ := net.ParseCIDR(r.Cidr)
		size, _ := mask.Mask.Size()
		addr3, addr4 := rand.IntN(255), rand.IntN(255)
		r.Cidr = fmt.Sprintf("192.168.%d.%d/%d", addr3, addr4, size)
		r.Name = strings.ToLower(fmt.Sprintf("%s-%s-%s", pokemons[rand.IntN(len(pokemons))], pokemons[rand.IntN(len(pokemons))], pokemons[rand.IntN(len(pokemons))]))
		return *r
	}
	return *r
}
