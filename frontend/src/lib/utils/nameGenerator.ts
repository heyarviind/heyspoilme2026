import { uniqueNamesGenerator, type Config } from 'unique-names-generator';

// Female first names - suggestive/playful
const femaleFirstNames: string[] = [
	'Velvet', 'Sugar', 'Blush', 'Honey', 'Peach', 'Pink', 'Silk', 'Rose', 'Candy', 'Cherry',
	'Vixen', 'Siren', 'Doll', 'Angel', 'Bunny', 'Desire', 'Tempt', 'Flame', 'Kiss', 'Crush',
	'Lux', 'Baby', 'Darling', 'Sweetie', 'Dream', 'Muse', 'Glow', 'Heat', 'Lush', 'Fever',
	'Bliss', 'Envy', 'Sin', 'Tease', 'Crave', 'Allure', 'Obsession', 'Fantasy', 'Eclipse', 'Voltage',
	'Spark', 'Bloom', 'Whisper', 'Pulse', 'Vibe', 'Magic', 'Charm', 'Desiree', 'Dollface', 'Seduce',
	'Kissy', 'Babe', 'Cutie', 'Naughty', 'Sassy', 'Minx', 'Kitty', 'Foxy', 'CherryPop', 'Sugarplum',
	'Blushie', 'Peachy', 'Honeybun', 'Sweetpea', 'Pinky', 'Silky', 'Rosy', 'Candyfloss', 'Dreamy', 'Softie',
	'Luna', 'Nova', 'Aura', 'Ember', 'Nyx', 'Zara', 'Roxy', 'Ivy', 'Cleo', 'Kali',
	'Jinx', 'Luxie', 'Vee', 'Elle', 'Mimi', 'Kiki', 'Coco', 'Gigi', 'Lola', 'Fifi',
	'Poppy', 'Skye', 'Raven', 'Sienna', 'Blaze', 'Misty', 'Stormy', 'Cherryx', 'Pinkie', 'Velvetx'
];

// Female second names
const femaleSecondNames: string[] = [
	'Kiss', 'Desire', 'Temptation', 'Heat', 'Sin', 'Crush', 'Flame', 'Allure', 'Fantasy', 'Tease',
	'Doll', 'Siren', 'Vixen', 'Muse', 'Babe', 'Angel', 'Bunny', 'Minx', 'Fox', 'Kitty',
	'Touch', 'Whisper', 'Spell', 'Obsession', 'Lush', 'Fever', 'Bliss', 'Rush', 'Spark', 'Glow',
	'Velvet', 'Silk', 'Sugar', 'Honey', 'Candy', 'Cherry', 'Peach', 'Pink', 'Rose', 'Bloom',
	'Pulse', 'Vibe', 'Magic', 'Charm', 'Dream', 'Eclipse', 'Voltage', 'Aura', 'Heatwave', 'Mirage',
	'Kissed', 'Desired', 'Tempted', 'Craved', 'Wanted', 'Loved', 'Obsessed', 'Hooked', 'Addicted', 'Taken',
	'Darling', 'Sweetheart', 'Lover', 'Baby', 'Dollface', 'Seductress', 'Goddess', 'Queen', 'Princess', 'Empress',
	'Naughty', 'Wild', 'Soft', 'Hot', 'Sweet', 'Dangerous', 'Private', 'Secret', 'Hidden', 'Forbidden',
	'Night', 'Midnight', 'Afterdark', 'Desiree', 'Pleasure', 'Ecstasy', 'Touché', 'Luxe', 'Divine', 'Fatal',
	'Fix', 'Addiction', 'Hunger', 'Craving', 'Fixation', 'Affair', 'Fantasyx', 'Kissme', 'Yours', 'Obey'
];

// Male first names - dominant/powerful
const maleFirstNames: string[] = [
	'Alpha', 'Apex', 'Prime', 'Sovereign', 'Monarch', 'Titan', 'Regent', 'Noble', 'Magnus', 'Victor',
	'Sterling', 'Valor', 'Baron', 'Duke', 'Kaiser', 'Caesar', 'Emperor', 'Lord', 'Knight', 'Paladin',
	'Atlas', 'Orion', 'Ares', 'Odin', 'Zeus', 'Thor', 'Mars', 'Hades', 'Apollo', 'Achilles',
	'Wolf', 'Lion', 'Panther', 'Hawk', 'Falcon', 'Raven', 'Cobra', 'Viper', 'Dragon', 'Bull',
	'Iron', 'Steel', 'Stone', 'Slate', 'Onyx', 'Obsidian', 'Flint', 'Granite', 'Titanium', 'Carbon',
	'Kingpin', 'Boss', 'Don', 'Chief', 'Captain', 'Marshal', 'Warden', 'Sentinel', 'Enforcer', 'Operator',
	'Phantom', 'Shadow', 'Specter', 'Ghost', 'Reaper', 'Hunter', 'Predator', 'Sniper', 'Rogue', 'Outlaw',
	'Blade', 'Edge', 'Fang', 'Claw', 'Hammer', 'Anvil', 'Saber', 'Dagger', 'Arrow', 'Cross',
	'Cipher', 'Vector', 'Nexus', 'Flux', 'Signal', 'Protocol', 'Legacy', 'Dynasty', 'Empire', 'Ascend',
	'Void', 'Night', 'Storm', 'Thunder', 'Frost', 'Inferno', 'Vortex', 'Eclipse', 'Nova', 'Zenith'
];

// Male second names
const maleSecondNames: string[] = [
	'Prime', 'Elite', 'Supreme', 'Absolute', 'Legacy', 'Dynasty', 'Empire', 'Dominion', 'Ascend', 'Apex',
	'Sovereign', 'Monarch', 'Regent', 'Noble', 'Imperial', 'Royal', 'Majestic', 'Crown', 'Throne', 'Prestige',
	'Valor', 'Honor', 'Creed', 'Code', 'Oath', 'Ethos', 'Virtue', 'Resolve', 'Command', 'Authority',
	'Steel', 'Iron', 'Stone', 'Onyx', 'Obsidian', 'Granite', 'Titanium', 'Carbon', 'Forge', 'Anvil',
	'Wolf', 'Lion', 'Hawk', 'Falcon', 'Panther', 'Dragon', 'Bull', 'Raven', 'Viper', 'Cobra',
	'Shadow', 'Phantom', 'Specter', 'Ghost', 'Night', 'Eclipse', 'Vortex', 'Storm', 'Thunder', 'Inferno',
	'Blade', 'Saber', 'Hammer', 'Dagger', 'Arrow', 'Edge', 'Fang', 'Claw', 'Cross', 'Shield',
	'Sentinel', 'Warden', 'Marshal', 'Enforcer', 'Guardian', 'Operator', 'Commander', 'Admiral', 'Captain', 'Chief',
	'Nexus', 'Vector', 'Cipher', 'Protocol', 'Signal', 'Core', 'Axis', 'Vertex', 'Primeval', 'Omega',
	'Fix', 'Drive', 'Force', 'Will', 'Instinct', 'Power', 'Focus', 'Control', 'DominionX', 'LegacyX'
];

export function generateDisplayName(gender: 'male' | 'female'): string {
	const firstNames = gender === 'male' ? maleFirstNames : femaleFirstNames;
	const secondNames = gender === 'male' ? maleSecondNames : femaleSecondNames;

	const config: Config = {
		dictionaries: [firstNames, secondNames],
		separator: ' ',
		length: 2,
		style: 'capital'
	};

	return uniqueNamesGenerator(config);
}

