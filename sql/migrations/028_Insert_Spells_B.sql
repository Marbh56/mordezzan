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
        'Barkskin',
        3,
        'wch,drd',
        'touch',
        '1 turn + 1 turn per CA level',
        'Toughens the recipient''s flesh to the strength of oak bark, providing an armour class equivalent to chain mail (AC 5, DR 1), or +1 AC if other armour is worn. Barkskin does not encumber the recipient in any way.'
    ),
    (
        'Befriend Animals',
        2,
        'wch,drd',
        '10 feet',
        'permanent',
        'The sorcerer enchants normal animals. The subjects must be Neutral, non-magical creatures of ordinary sort: amphibians, birds, fish, mammals, or reptiles. The sorcerer can control 1 HD of animals per CA level. Each such animal is granted a sorcery saving throw. Torpid or non-hostile animals make normal saving throws, but attacking animals save at +2. Animals that fail their saves are mesmerized and will follow the sorcerer to the best of their ability; those that make their saving throws will wander away or attack the caster, depending on the animal''s nature. An enchanted animal can be taught three tricks on par with what a trained dog or cat might learn, depending on its intelligence and capability. Each trick requires about a week of uninterrupted training; i.e., not during the course of adventure.'
    ),
    (
        'Black Cloud',
        3,
        'mag,cry,wch,drd',
        '240 feet',
        '1 turn',
        'A small raincloud appears 30–60 feet above the intended target area. It releases a torrent of rain that forms a cylinder, showering a 30-foot-diameter area. If this spell is cast in an area of subfreezing temperature, the precipitation instead will be heavy snow; or, if the temperature is just at the freezing point, sleet and freezing rain results. All attack rolls made whilst under a black cloud are at −4 "to hit" penalties. Normal fires will be extinguished; magical fires will be temporarily snuffed, their dweomers rekindling 1 turn after the spell terminates (unless their durations elapse). Black cloud can also be used as a protective measure, for if a fireball, flaming sphere, or similar effect strikes the deluged area, the fire spell will be extinguished, and the rain vaporized to a cloud of steam.'
    ),
    (
        'Black Hand',
        1,
        'nec',
        '0',
        '3 rounds + 1 round per CA level',
        'The sorcerer''s right hand turns dark as pitch and emits tiny motes of black and violet. The black hand enhances touch spells of harmful intent, such as inflict disease, ghoul touch, and shocking grasp. These subsequent touch attacks are made at +1 "to hit" for every 4 CA levels: CA 1–4 = +1, CA 5–8 = +2, CA 9–12 = +3. The black hand spell does not expire after a touch spell is successfully delivered; it persists for the full duration noted above.'
    ),
    (
        'Black Tentacles',
        4,
        'mag,nec,wch',
        '30 feet',
        '1 round per CA level',
        'Black, squid-like tentacles erupt in a 30-foot-diameter area, one such tentacle per CA level of the sorcerer. Each thick and slimy tentacle is 10 feet long, AC 4, and equal in hit points to the sorcerer at full health. Any creature within range of a black tentacle is subject to attack; if a tentacle has more than one potential target, the referee should assign equal chances via random die roll. Each victim must make a sorcery saving throw. If the saving throw succeeds, the tentacle lashes the target for 1d6 hp damage before disappearing. If the saving throw fails, the tentacle lash delivers 1d6 hp damage, as well as constricting and rending the victim for a further 2d6 hp damage per round until the spell ends or the tentacle is destroyed. As the tentacles have no intelligence, they will continue to squeeze a dead body and might on occasion be fooled into constricting a barrel, statue, tree, or like object.'
    ),
    (
        'Blade Barrier',
        6,
        'clr',
        '30 feet',
        '1 turn',
        'The sorcerer conjures a 12-foot wall of whirling, keen-edged blades that spin and flash around a selected point, fencing in an area as small as 5 × 5 feet to as large as 50 × 50 feet. Any creature that attempts to pass through the blade barrier will be assailed by the whirling blades, sustaining 8d8 hp damage. When this spell is cast, targeted creatures are entitled to an avoidance saving throw to escape harm; however, there is a 3-in-6 chance that they escape within the blade barrier, not without.'
    ),
    (
        'Bless',
        2,
        'wch,clr',
        '0',
        '3 turns',
        'All allies within 25 feet of the caster are sanctified by this spell, gaining a +1 to any saving throw versus fear effects (whether sorcery or device), and a +1 bonus on all attack rolls. Furthermore, NPCs (henchmen, hirelings, etc.) each gain a +1 bonus on any morale (ML) check. The reverse form of this spell, blight, curses all hostile creatures within 25 feet of the caster, effecting −1 morale, −1 on saving throws versus fear effects, and a −1 penalty on all attack rolls. N.B.: Either form of this spell affects only those within range at the moment the spell is cast; i.e., subsequently moving into or out of range has no bearing on the spell''s effects.'
    ),
    (
        'Breathe Fire',
        5,
        'pyr,drd',
        '10 feet',
        'special',
        'The lips of the sorcerer must be pursed after speaking the final incantation of this spell, for the next time the sorcerer''s mouth is opened, a jet of flames 10 feet long and 5 feet wide at its terminus is released. Victims in this path sustain 3d8+3 hp damage, though they can attempt avoidance saving throws for half damage. The sorcerer''s mouth may be opened at will to release this spell, perhaps during combat or other like activities; however, other spells may not be cast. If breathe fire is not released within 1 turn (10 minutes), the sorcerer immolates, sustaining maximum damage (27 hp) with no saving throw applicable. (This spell can be dangerous if the caster is forgetful and speaks to an ally or another person.)'
    ),
    (
        'Breathe Frost',
        5,
        'cry',
        '10 feet',
        'special',
        'The lips of the sorcerer must be pursed after speaking the final incantation of this spell, for the next time the sorcerer''s mouth is opened, a billowing jet of frost 10 feet long and 5 feet wide at its terminus is released. Victims in this path sustain 3d8+3 hp damage, though they can attempt avoidance saving throws for half damage. The sorcerer''s mouth may be opened at will to release this spell, perhaps during combat or other like activities; however, other spells may not be cast. If breathe frost is not released within 1 turn (10 minutes), the sorcerer internally freezes, sustaining maximum damage (27 hp) with no saving throw applicable. (This spell can be dangerous if the caster is forgetful and speaks to an ally or another person.)'
    ),
    (
        'Brink of Death',
        5,
        'nec,clr',
        'touch',
        'instantaneous',
        'Revives a just-killed human or other creature, providing the spell is cast within 6 rounds (1 minute) of expiry. The subject must make a trauma survival check (see Chapter 3: Statistics, constitution) and furthermore suffers a permanent loss of 1 point of constitution. Brink of death also can be used to bring back a living but unconscious subject from a negative hit point total. The subject is immediately restored to consciousness (at 1 hp). Casting the spell in this manner entails neither a trauma survival check nor constitution loss.'
    ),
    (
        'Burning Hands',
        1,
        'mag,pyr',
        '5 feet',
        'instantaneous',
        'Jets of thin, multihued flames spring from the fingertips of the caster''s enveloped-in-flames hands, fanning out in a 120° horizontal arc and causing 2 hp damage per CA level, with no saving throw allowed. Combustible materials (e.g., cloth, paper, dry wood) are likely ignited if exposed to burning hands.'
    );

-- +goose Down
DELETE FROM spells
WHERE
    name IN (
        'Barkskin',
        'Befriend Animals',
        'Black Cloud',
        'Black Hand',
        'Black Tentacles',
        'Blade Barrier',
        'Bless',
        'Breathe Fire',
        'Breathe Frost',
        'Brink of Death',
        'Burning Hands'
    );
