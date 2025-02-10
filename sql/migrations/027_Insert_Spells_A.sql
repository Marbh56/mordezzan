-- +goose Up
INSERT INTO
    spells (
        name,
        level,
        classes,
        range,
        duration,
        description
    )
VALUES
    (
        'Acid Arrow',
        2,
        'mag',
        '30 feet',
        'special',
        'A magic arrow darts from the finger of the caster. On a successful attack roll (dexterity modifier applies), the acid arrow strikes for 1d4+1 hp physical damage, plus an additional 2d4 hp acid damage in the same round. Magicians (but not other sorcerers) enjoy a +1 bonus to the attack roll for every 2 CA levels (CA 3–4 = +2, CA 5–6 = +3, and so forth). Acid damage will persist for higher level sorcerers: 1 extra round for CA 4–6, 2 extra rounds for CA 7–9, 3 extra rounds for CA 10 or greater. For example, an acid arrow fired by a CA 12 sorcerer on round 1 would inflict 1d4+1 hp base damage plus 2d4 hp acid damage on round 1, 2d4 hp acid damage on round 2, 2d4 hp acid damage on round 3, and a final 2d4 hp acid damage on round 4. The acid may ruin armour or clothing per referee discretion. N.B.: If using the optional critical hits and misses rules, acid arrow is subject to critical success or failure; however, magicians should use the "fighter" column in each such instance. Also, any multiplied damage that results from a critical hit or critical miss applies strictly to the arrow''s physical damage, not the acid.'
    ),
    (
        'Acid Fog',
        6,
        'ill',
        '30 feet',
        '6 rounds (1 minute)',
        'Creates a caustic miasma as large as 10 feet thick, 30 feet long, and 30 feet high. Only 35 mph or greater wind will disperse acid fog; even a gust of wind spell is ineffective. Acid fog burns all within its confines and is particularly deadly to vegetation, killing small plants and grass at once. Humans, monsters, and other fauna are harmed as well, the acid blistering their skin, eyes, throat, and lungs. Such creatures sustain 1 hp damage on the 1st round, 2 hp on the 2nd, 4 hp on the 3rd, 8 hp on the 4th, 10 hp on the 5th, and 12 hp on the 6th; vegetal creatures suffer 150% of these damage figures. Any creature caught within or attempting to traverse acid fog is slowed by the fog''s opacity to a maximum rate of 10 MV. Normal sight cannot penetrate acid fog, and infrared vision is unavailing. Potent fire spells such as fireball, flame strike, or wall of fire will burn off acid fog in 1 round.'
    ),
    (
        'Advanced Hypnotism',
        5,
        'ill',
        '30 feet',
        'special',
        'To hypnotize a single target. To resist, the target is granted a sorcery saving throw, modified by willpower adjustment, if applicable. Hostile or aggressive creatures gain +1 to +3 bonuses on their saves, as judged by the referee. A failed save indicates the victim is unaware of being ensorcelled and is subject to a post hypnotic suggestion. The trigger is determined by the sorcerer; it may be something seen or a spoken word. The player must inform the referee what the specific trigger is, and the referee must judge if it is appropriate. Once triggered, a suggestion (q.v.) is effected. The spell is broken after this suggestion is triggered and acted on, which might be hours, days, months, or even years later.'
    ),
    (
        'Advanced Spectral Phantasm',
        5,
        'ill',
        '240 feet',
        'permanent',
        'A visual illusion is created, a projected image of nearly anything the caster can imagine, as large as 50 × 50 × 50 feet in area. Victims must be living creatures of animal intelligence or greater; undead, constructs, oozes, and the like are unaffected. Advanced spectral phantasm can be used to create an attacking monster or another damaging hazard. The illusion can be of sight, sound, smell, and/or temperature. The auditory component might include a shout, a roar, or a sentence of no more than nine words in length (not including articles a, an, and the). Once cast, this spell will persist infinitely, with no need of concentration unless the caster wishes to manipulate the movements of the illusion; such manoeuvres require full concentration, though the sorcerer can walk at half speed and maintain control. The illusion is broken if struck for 1 hp damage or more, or it can be terminated via a dispel phantasm spell. An advanced spectral phantasmal monster is AC 6 and will disappear if hit; otherwise, it can continue to attack without caster direction. Advanced spectral phantasm targets are not allowed saving throws unless the referee feels the illusion is unbelievable to its viewers, in which case sorcery saving throws should be rolled, modified by willpower adjustment, if applicable. With sight, sound, smell, and/or temperature, credibility is rarely an issue with this spell; if such a case arises, and the save is made, the disbeliever will see the advanced spectral phantasm as a flawed and flickering transparent image. An intelligent disbeliever may then alert allies, whose saves are made at a bonus of +4. Otherwise, this illusion can inflict real physical damage: 1d8 hp per CA level per round to each victim. Such damage remains even if the spell is subsequently broken.'
    ),
    (
        'Agonizing Touch',
        3,
        'nec',
        'touch',
        'instantaneous',
        'When touched by the sorcerer, a searing pain runs through the victim''s nervous system; no saving throw applies. Damage inflicted is 1d4 hp, but the intense pain impairs the victim for 1d4+1 rounds: The victim''s attack rolls, armour class, and saving throws are each at a −2 penalty; movement is halved; and spells or innate sorcerous abilities stand a 3-in-6 chance of failure (the optional concentration rule does not apply).'
    ),
    (
        'Aid',
        2,
        'clr',
        'touch',
        '1 turn',
        'A single recipient gains a +1 bonus on any saving throw versus fear effects (whether sorcery or device), a +1 bonus on all attack rolls, and a temporary 1d8 hit point boost. These hit points may exceed the recipient''s normal hp maximum. Any subsequent damage is drawn from the temporary hp first, the remainder disappearing when the spell ends. Lastly, if cast upon an NPC ally (henchman, hireling, etc.), the recipient gains a +1 bonus on any morale (ML) check.'
    ),
    (
        'Air Walk',
        5,
        'clr',
        'touch',
        '6 turns (1 hour)',
        'The sorcerer can walk on air as though it were solid ground and lead up to one recipient per CA level to do the same. Each recipient must be touched and be willing. In a straight, single-file line (such as over a chasm, ravine, or trench), the sorcerer can lead recipients to walk (but not run) at normal movement rate. Air walkers also can walk up or down at a 45° angle, as though ascending or descending stairs, at one-half normal walking speed. Too, they can ascend or descend vertically, as though climbing a sheer cliff with ample handholds and toeholds, at one-fourth normal walking speed. Any recipients that stray from the caster''s path will fall and take damage.'
    ),
    (
        'Air-like Water',
        5,
        'mag',
        '0',
        '6 turns (1 hour) per CA level',
        'Transforms a 30-foot-diameter sphere of fresh or salt water into a magical, bubbling solution that can be inhaled safely by air-breathing creatures; furthermore, underwater pressure is negated. The spell can be cast as the sorcerer enters the water or after submerging, the air-like water moving with the caster. Water-breathing creatures will instinctively avoid the sphere, as they cannot respire within its confines, but this general eschewal does not preclude certain predators of the deep from attempting to snatch prey from the air-like water.'
    ),
    (
        'Alarm',
        1,
        'mag',
        '10 feet per CA level',
        '12 turns (2 hours) per CA level',
        'Cast upon as many doors, gates, portals, or other point of ingress/egress that are within the caster''s range. This spell is triggered by the passage of any living creature larger than a rat (3+ lbs.), evoking a sound not unlike bells pealing. Undead, constructs, and other non-living entities will not activate the alarm spell; neither will incorporeal beings, though invisible creatures will set off the spell.'
    ),
    (
        'Allay Exhaustion',
        3,
        'ill',
        'touch',
        '6 turns (1 hour)',
        'Creates the illusion of healing, wellness, energy, and stamina. Allay exhaustion allows one to persevere without sleep when thoroughly exhausted, as though an extraordinary feat of constitution had been achieved. As well, any previous hit point loss is temporarily healed by 50%. The exact number should be recorded, for when the spell''s duration elapses, this illusory hit point boon is lost. Once the illusion ends, the recipient must rest for 12 turns (2 hours) or suffer −4 penalties to attack rolls, damage rolls, and saving throws. Unwilling recipients are allowed sorcery saving throws, modified by willpower adjustment, if applicable.'
    ),
    (
        'Alter Self',
        2,
        'ill,wch',
        '0',
        '1d6 turns',
        'The sorcerer''s form is mutated to something or someone no more than 50% smaller or larger, lighter or heavier. The new form has quasi-actuality. For example, if the new form has wings, the caster is allowed minimal flight, moving at no greater than 50% of the actual creature''s speed; if the creature has gills, the caster can breathe underwater; and so on. The altered form is limited to humans, humanoids, or other bipedal species with which the caster is familiar. This spell does not grant any special abilities beyond locomotion and respiration—no innate magical abilities, no enhanced acuity. The duration of the spell should be rolled in secret by the referee.'
    ),
    (
        'Animal Growth',
        5,
        'drd',
        '120 feet',
        '6 rounds (1 minute) per CA level',
        'Causes as many as six normal beasts (amphibians, birds, fish, mammals, or reptiles; not humans, humanoids, or monsters) to double in size. The effect results in doubled hit dice, doubled damage dice, and whatever else the referee deems appropriate. The reverse of this spell, animal reduction, shrinks as many as six animals to half their normal size, resulting in halved hit dice and halved damage on attacks. No saving throw is permitted for either form of this spell.'
    ),
    (
        'Animate Carrion',
        1,
        'nec',
        '10 feet',
        'permanent',
        'Raised are the bones or carrion of Small animals: amphibians, birds, mammals, and reptiles of natural sort. The small undead animals will obey the simple instructions of the caster (essentially one-word commands) and follow the sorcerer unless slain or turned; the dispel magic spell also nullifies the connexion between the sorcerer and the undead animal. The caster can animate and maintain no more than 1 HD of undead animals per CA level. Animated carrion loses any special abilities possessed in life (e.g., flight, musk, venom).'
    ),
    (
        'Animate Carrion II',
        3,
        'nec',
        '10 feet',
        'permanent',
        'Raised are the bones or carrion of Medium animals: amphibians, birds, mammals, and reptiles of natural sort. The medium undead animals will obey the simple instructions of the caster (essentially one-word commands) and follow the sorcerer unless slain or turned; the dispel magic spell also nullifies the connexion betwixt the sorcerer and the undead animal. The caster can animate and maintain no more than 2 HD of undead animals per CA level. Animated carrion loses any special abilities possessed in life (e.g., flight, musk, venom).'
    ),
    (
        'Animate Carrion III',
        5,
        'nec',
        '10 feet',
        'permanent',
        'Raised are the bones or carrion of Large animals: amphibians, birds, mammals, and reptiles of natural sort. The large undead animals will obey the simple instructions of the caster (essentially one-word commands) and follow the sorcerer unless slain or turned; the dispel magic spell also nullifies the connexion betwixt the sorcerer and the undead animal. The caster can animate and maintain no more than 3 HD of undead animals per CA level. Animated carrion loses any special abilities possessed in life (e.g., flight, musk, venom).'
    ),
    (
        'Animate Dead',
        5,
        'mag,nec,wch,clr',
        '10 feet',
        'permanent',
        'Skeletons or zombies are created from the bones and cadavers of dead humans or humanoids. The undead will obey the commands of the caster, following, attacking, or standing guard as directed. They will continue to serve until slain or turned; the dispel magic spell also nullifies the connexion betwixt the sorcerer and the undead. Through this necromancy the sorcerer can animate and maintain 1 skeleton or zombie per CA level. If suitable remains are at hand, the sorcerer can opt to raise 1 large skeleton per 3 CA levels, or 1 giant skeleton per 6 CA levels, though zombies may only be created from the whole corpses of humans.'
    ),
    (
        'Animate Dead II',
        6,
        'nec',
        '10 feet',
        'permanent',
        'Unspeakable rites and forbidden incantations raise ghouls from the fresh graves of humans. The selected graves must be no older than one week and dug properly. The ghouls will claw out from the earth to obey the commands of the sorcerer, following, attacking, or standing guard as directed. They will continue to serve until either slain or turned; the dispel magic spell also nullifies the connexion betwixt the sorcerer and the undead. Through this necromancy the sorcerer can animate and maintain 1 ghoul for every 2 CA levels. If a 12th-level sorcerer casts this spell, the sixth ghoul will emerge as a ghast.'
    ),
    (
        'Animate Objects',
        6,
        'wch,clr',
        '60 feet',
        '6 turns (1 hour)',
        'Enchants non-magical articles to rouse and do the sorcerer''s bidding, affecting one or more objects of total weight not exceeding 400 pounds. The referee should determine movement rate, hit points, attacks, and damage delivered by the animated objects. Consider the following guidelines: Boulder, Round (250 lbs.): MV 20; DX 5; AC 4; HD 5; #A 1/1 (rolling smash); D 1d10 Chest, Iron: MV 10; DX 4; AC 5; HD 3; #A 1/1 (bite); D 1d8 Statue, Stone: MV 20; DX 6; AC 1; HD 6; #A 1/1 (strike); D 2d8 Table, Wooden: MV 30; DX 4; AC 7; HD 2; #A 2/1 (legs); D 1d6/1d6 Animated objects use the fighting ability (FA) of the caster and make saving throws according to their item category; morale (ML) does not apply.'
    ),
    (
        'Anti-Beast Shell',
        6,
        'drd',
        '0',
        '1 turn per CA level',
        'Creates an invisible, 10-foot-radius, hemispherical field around the sorcerer. The barrier prevents any animal from breaching the anti-beast shell or attacking those within. This spell does not affect magical beasts or "monsters"; only natural beasts of the animal kingdom are hedged out: amphibians, arachnids, birds, fish, insects, mammals, and reptiles, including giant-sized species. Those afforded the protection of this spell cannot attack or otherwise harm any animal outside the shell, or the spell will terminate.'
    ),
    (
        'Anti-Magic Field',
        6,
        'mag,wch',
        '0',
        '12 rounds (2 minutes)',
        'Evokes a magical energy shield to surround the sorcerer at a radius of one foot per CA level. The anti-magic field repels any spell or sorcerous effect (as from a ring, staff, wand, etc.); however, just as no spell or spell effect can enter the shell, no spell or spell effect can exit it, either.'
    ),
    (
        'Anti-Plant Shell',
        5,
        'drd',
        '0',
        '1 turn per CA level',
        'Creates an invisible, 10-foot-radius, hemispherical field that encircles the sorcerer. The barrier prevents any plant from breaching the anti-plant shell or attacking him, including vegetal monsters such as green slime, mustard mould, shambling mounds, and tree-men. Those afforded the protection of this spell cannot attack or otherwise harm any plant creature outside the shell, or the spell will terminate.'
    ),
    (
        'Atonement',
        5,
        'clr',
        'touch',
        'permanent',
        'Usually cast on those of similar religion and/or like alignment, this spell takes 1 hour to cast, following prayer, cogitation, and incense burning. It removes the onus of misdeeds that are unknowingly, unintentionally, or unwillingly committed; also, this spell can undo the effects of magical alignment change. If the recipient has exercised poor judgment and consequently violated the precepts of faith and/or alignment, this spell can remove the burden or penalties accorded if the character is truly repentant. Ultimately, the subject''s contrition must be judged by the referee. One cannot atone for deliberate misdeeds. The recipient of the atonement spell might be charged with a quest (q.v.) to complete his reparations.'
    ),
    (
        'Auditory Glamour',
        2,
        'mag,ill',
        '240 feet',
        '1 turn',
        'A hallucination of sound is created, that of voices, calls, or cries (human, humanoid, animal, or monster); footfalls; or other like noises. The sounds of 1d4 creatures can be invented thus for each CA level of the sorcerer. However, if a sound is of significant volume, the referee must decide on the number of voices and their collective volume (e.g., the roar of one lion may be equal to the shouts of five people).'
    ),
    (
        'Augury',
        2,
        'clr',
        '0',
        'special',
        'Through communion with otherworldly agents, the sorcerer learns whether an action in the near future (within 3 turns) will be advantageous or disadvantageous. The caster must clearly and concisely articulate the considered action (or inaction, as it were) through prayer and cogitation over a 1-turn period. The referee then informs the player if the proposed action is for weal, for woe, or inconsequential. The chance to successfully divine the future is 7-in-10, which should be rolled secretly by the referee; a failed result yields an inaccurate augury.'
    );

-- +goose Down
DELETE FROM spells
WHERE
    name IN (
        'Acid Arrow',
        'Acid Fog',
        'Advanced Hypnotism',
        'Advanced Spectral Phantasm',
        'Agonizing Touch',
        'Aid',
        'Air Walk',
        'Air-like Water',
        'Alarm',
        'Allay Exhaustion',
        'Alter Self',
        'Animal Growth',
        'Animate Carrion',
        'Animate Carrion II',
        'Animate Carrion III',
        'Animate Dead',
        'Animate Dead II',
        'Animate Objects',
        'Anti-Beast Shell',
        'Anti-Magic Field',
        'Anti-Plant Shell',
        'Atonement',
        'Auditory Glamour',
        'Augury'
    );
