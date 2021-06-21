package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/particles"
)

func UseLaserGun(info ItemUseInfo) {
	// click relative to camera
	camCoords :=
		mgl32.Vec4{
			info.ScreenX / render.PixelsPerTile,
			info.ScreenY / render.PixelsPerTile, 0, 1 }
	// click relative to world
	worldCoords := info.Camera.GetTransform().Mul4x1(camCoords)
	// adding 0.5 here because raytrace assumes top left coords,
	// but player and tile positions are stored using centered coords
	positions := cxmath.Raytrace(
		float64(info.Player.Pos.X)+0.5,float64(info.Player.Pos.Y)+0.5,
		float64(worldCoords.X())+0.5,float64(worldCoords.Y()) + 0.5 )

	playerPos := mgl32.Vec2{ info.Player.Pos.X, info.Player.Pos.Y}
	for _,pos := range positions {
		if info.Planet.TileIsSolid(int(pos.X),int(pos.Y)) {
			direction := worldCoords.Vec2().Sub(playerPos).Normalize()
			length := pos.Vec2().Sub(playerPos).Len() + 0.5
			targetPos := playerPos.Add(direction.Mul(length))

			particles.CreateLaser(playerPos, targetPos)

			info.Planet.DamageTile(int(pos.X), int(pos.Y), world.TopLayer)
			return
		}
	}

	// hit nothing - visual effect only
	particles.CreateLaser(playerPos, worldCoords.Vec2() )
}

func RegisterLaserGunItemType() ItemTypeID {
	// TODO use proper asset
	laserGunSpriteId :=spriteloader.LoadSingleSprite(
			"./assets/item/lasergun-temp.png","lasergun")
	laserGunItemType := NewItemType(laserGunSpriteId)
	laserGunItemType.Use = UseLaserGun
	laserGunItemType.Name = "Laser Gun"
	return AddItemType(laserGunItemType)
}
