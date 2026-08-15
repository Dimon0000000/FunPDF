export interface BaseModel {
  create_time?: number
  create_date?: string
  update_time?: number
  update_date?: string
}

export interface Album extends BaseModel {
  id: string
  name: string
  avatar: string
  description: string
  files?: AlbumFile[]
}

export interface AlbumFile extends BaseModel {
  id: string
  name: string
  album_id?: string
  sort_order?: number
  description?: string
}

export interface Provider extends BaseModel {
  id: string
  name: string
  enabled?: boolean
  base_url?: string
  models?: ProviderModel[]
}

export interface ProviderModel extends BaseModel {
  id: string
  name: string
  provider_name: string
}
