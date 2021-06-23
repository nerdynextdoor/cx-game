package world

import (
	"math/rand"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/spriteloader"
)

// TODO shove in .yaml file
const persistence = 0.5
const lacunarity = 2

func (planet *Planet) placeTileOnTop(x int, tile Tile) {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers.Top[tileIdx] = tile
}

func (planet *Planet) placeLayer(tile Tile, depth,noiseScale float32) {
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), 4, 256)
	for x:=int32(0); x<planet.Width; x++ {
		noiseSample := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		height := int(depth+noiseSample*noiseScale)
		for i:=0; i<height; i++ {
			planet.placeTileOnTop(int(x),tile)
		}
	}
}


func GeneratePlanet() *Planet {
	planet := NewPlanet(100, 100)
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png", "Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")

	dirt := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
		Name: "Dirt",
	}
	stone := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Stone")),
		Name: "Stone",
	}
	bedrock := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Bedrock")),
		Name: "Bedrock",
	}

	planet.placeLayer(bedrock, 4,1)
	planet.placeLayer(stone, 4,1)
	planet.placeLayer(dirt, 4,1)
	
	return planet
}
