package core

import "strings"

// Field is a stable semantic identifier. Identifiers are serialized as strings.
type Field string

const (
	Title                       Field = "title"
	Subtitle                    Field = "subtitle"
	Album                       Field = "album"
	Grouping                    Field = "grouping"
	Version                     Field = "version"
	OriginalTitle               Field = "original_title"
	Artist                      Field = "artist"
	ArtistDisplay               Field = "artist_display"
	Artists                     Field = "artists"
	AlbumArtist                 Field = "album_artist"
	AlbumArtistDisplay          Field = "album_artist_display"
	OriginalArtist              Field = "original_artist"
	Performer                   Field = "performer"
	FeaturedArtist              Field = "featured_artist"
	PerformerCredits            Field = "performer_credits"
	Composer                    Field = "composer"
	Lyricist                    Field = "lyricist"
	Conductor                   Field = "conductor"
	Arranger                    Field = "arranger"
	Work                        Field = "work"
	Movement                    Field = "movement"
	MovementNumber              Field = "movement_number"
	MovementTotal               Field = "movement_total"
	ShowMovement                Field = "show_movement"
	OpusNumber                  Field = "opus_number"
	Producer                    Field = "producer"
	Engineer                    Field = "engineer"
	Mixer                       Field = "mixer"
	Remixer                     Field = "remixer"
	MasteringEngineer           Field = "mastering_engineer"
	ProductionCredits           Field = "production_credits"
	Publisher                   Field = "publisher"
	RecordLabel                 Field = "record_label"
	Organization                Field = "organization"
	CatalogNumber               Field = "catalog_number"
	TrackNumber                 Field = "track_number"
	TrackTotal                  Field = "track_total"
	DiscNumber                  Field = "disc_number"
	DiscTotal                   Field = "disc_total"
	Compilation                 Field = "compilation"
	MediaType                   Field = "media_type"
	MediaKind                   Field = "media_kind"
	Genre                       Field = "genre"
	Mood                        Field = "mood"
	BPM                         Field = "bpm"
	MusicalKey                  Field = "musical_key"
	Language                    Field = "language"
	RecordingDate               Field = "recording_date"
	ReleaseDate                 Field = "release_date"
	OriginalReleaseDate         Field = "original_release_date"
	EncodingDate                Field = "encoding_date"
	TaggingDate                 Field = "tagging_date"
	Copyright                   Field = "copyright"
	License                     Field = "license"
	Ownership                   Field = "ownership"
	PaymentURL                  Field = "payment_url"
	TermsOfUse                  Field = "terms_of_use"
	ISRC                        Field = "isrc"
	UPC                         Field = "upc"
	Barcode                     Field = "barcode"
	ASIN                        Field = "asin"
	RecordingID                 Field = "musicbrainz_recording_id"
	ReleaseTrackID              Field = "musicbrainz_release_track_id"
	ReleaseID                   Field = "musicbrainz_release_id"
	ReleaseGroupID              Field = "musicbrainz_release_group_id"
	ArtistID                    Field = "musicbrainz_artist_id"
	AlbumArtistID               Field = "musicbrainz_album_artist_id"
	WorkID                      Field = "musicbrainz_work_id"
	DiscID                      Field = "musicbrainz_disc_id"
	AcoustID                    Field = "acoustid_id"
	AcoustIDFingerprint         Field = "acoustid_fingerprint"
	MusicIPFingerprint          Field = "musicip_fingerprint"
	BinaryIdentifier            Field = "binary_identifier"
	TitleSort                   Field = "title_sort"
	ArtistSort                  Field = "artist_sort"
	AlbumSort                   Field = "album_sort"
	AlbumArtistSort             Field = "album_artist_sort"
	ComposerSort                Field = "composer_sort"
	Comments                    Field = "comments"
	Description                 Field = "description"
	LongDescription             Field = "long_description"
	Lyrics                      Field = "lyrics"
	SynchronizedLyrics          Field = "synchronized_lyrics"
	Chapters                    Field = "chapters"
	TimedMetadata               Field = "timed_metadata"
	Podcast                     Field = "podcast"
	PodcastShow                 Field = "podcast_show"
	PodcastEpisode              Field = "podcast_episode"
	PodcastSeason               Field = "podcast_season"
	PodcastCategory             Field = "podcast_category"
	PodcastFeedURL              Field = "podcast_feed_url"
	PodcastEpisodeID            Field = "podcast_episode_id"
	Artwork                     Field = "artwork"
	ReplayGainTrackGain         Field = "replaygain_track_gain"
	ReplayGainAlbumGain         Field = "replaygain_album_gain"
	ReplayGainTrackPeak         Field = "replaygain_track_peak"
	ReplayGainAlbumPeak         Field = "replaygain_album_peak"
	ReplayGainReferenceLoudness Field = "replaygain_reference_loudness"
	R128TrackGain               Field = "r128_track_gain"
	R128AlbumGain               Field = "r128_album_gain"
	IntegratedLoudness          Field = "integrated_loudness"
	LoudnessRange               Field = "loudness_range"
	Rating                      Field = "rating"
	ContentRating               Field = "content_rating"
	PlayCount                   Field = "play_count"
	GaplessAlbum                Field = "gapless_album"
	OfficialArtistURL           Field = "official_artist_url"
	PublisherURL                Field = "publisher_url"
	PurchaseURL                 Field = "purchase_url"
	CopyrightURL                Field = "copyright_url"
	AudioSourceURL              Field = "audio_source_url"
	OfficialAudioURL            Field = "official_audio_url"
	Encoder                     Field = "encoder"
	EncodingSettings            Field = "encoding_settings"
	MediaOwner                  Field = "media_owner"
	Contact                     Field = "contact"
	Location                    Field = "location"
	ReleaseStatus               Field = "release_status"
	ReleaseType                 Field = "release_type"
	ReleaseCountry              Field = "release_country"
	Script                      Field = "script"
	DiscSubtitle                Field = "disc_subtitle"
	OriginalFilename            Field = "original_filename"
	OriginalLyricist            Field = "original_lyricist"
	PlaylistDelay               Field = "playlist_delay"
	RadioStationName            Field = "radio_station_name"
	RadioStationOwner           Field = "radio_station_owner"
	ISWC                        Field = "iswc"
	OriginalAlbum               Field = "original_album"
	InitialProducer             Field = "initial_producer"
	AudioDuration               Field = "audio_duration"
	EncoderDelay                Field = "encoder_delay"
	EncoderPadding              Field = "encoder_padding"
)

// Format is the metadata format and, for ID3, its on-disk major version.
type Format string

const (
	MP3v1  Format = "ID3v1"
	MP3v22 Format = "ID3v2.2"
	MP3v23 Format = "ID3v2.3"
	MP3v24 Format = "ID3v2.4"
	MP4    Format = "MP4"
	FLAC   Format = "FLAC"
	Vorbis Format = "Ogg Vorbis"
	Opus   Format = "Ogg Opus"
)

type Classification string

const (
	Native       Classification = "Native"
	Conventional Classification = "Conventional"
	Fallback     Classification = "Fallback"
	Partial      Classification = "Partial"
	ReadOnly     Classification = "Read-only"
	Unsupported  Classification = "Unsupported"
)

type Representation struct {
	Class     Classification
	Key       string
	Aliases   []string
	Note      string
	Reference string
}

type Definition struct {
	ID          Field
	Name        string
	Meaning     string
	Kind        Kind
	Multiple    bool
	Unit        string
	Constraints string
	Mutable     bool

	comment, id3v24, mp4, id3v23, id3v22 string
}

// Mapping describes the intended representation. A partial mapping requires
// format-specific limits or a lossless fallback for values it cannot express.
func (d Definition) Mapping(format Format) Representation {
	var key, ref string
	class := Native
	switch format {
	case MP3v1:
		key = map[Field]string{Title: "title", Artist: "artist", Album: "album", RecordingDate: "year", Comments: "comment", TrackNumber: "track (1.1)", Genre: "genre index"}[d.ID]
		ref = "https://id3.org/ID3v1"
		if key != "" {
			return Representation{Class: Partial, Key: key, Reference: ref, Note: "Legacy fixed-width input; canonical edits use ID3v2."}
		}
	case MP3v22:
		key, ref = d.id3v22, "https://id3.org/id3v2-00"
	case MP3v23:
		key, ref = d.id3v23, "https://id3.org/id3v2.3.0"
		if key == "" && (strings.HasPrefix(d.id3v24, "~TXXX:") || strings.HasPrefix(d.id3v24, "~WXXX:")) {
			key = d.id3v24
		}
	case MP3v24:
		key, ref = d.id3v24, "https://id3.org/id3v2.4.0-frames"
	case MP4:
		key, ref = d.mp4, "https://mutagen.readthedocs.io/en/latest/api/mp4.html"
		class = Conventional
	case FLAC, Vorbis, Opus:
		key, ref = d.comment, "https://picard-docs.musicbrainz.org/en/latest/appendices/tag_mapping.html"
		if strings.HasPrefix(key, "^") {
			return Representation{Class: Fallback, Key: CommentFallbackName(d.ID), Aliases: strings.Split(key[1:], ","), Note: "Recognizes compatibility input; no cited cross-player standard for writing it."}
		}
		class = Conventional
		if d.ID == AudioDuration && format == FLAC {
			return Representation{Class: ReadOnly, Key: "STREAMINFO", Reference: "https://datatracker.ietf.org/doc/html/rfc9639#section-8.2"}
		}
		if d.ID == Artwork && format == FLAC {
			return Representation{Class: Native, Key: "PICTURE (block 6)", Reference: "https://datatracker.ietf.org/doc/html/rfc9639#section-8.8"}
		}
		if d.ID == Artwork {
			ref = "https://wiki.xiph.org/VorbisComment#Unofficial_Field_names"
		}
		if format == Opus {
			switch d.ID {
			case R128TrackGain, R128AlbumGain:
				return Representation{Class: Native, Key: strings.ToUpper(string(d.ID)), Reference: "https://datatracker.ietf.org/doc/html/rfc7845#section-5.2.1"}
			case ReplayGainTrackGain, ReplayGainAlbumGain, ReplayGainTrackPeak, ReplayGainAlbumPeak, ReplayGainReferenceLoudness:
				return Representation{Class: Fallback, Key: CommentFallbackName(d.ID), Aliases: []string{key}, Note: "RFC 7845 discourages ReplayGain keys for Opus. R128 is a separate field.", Reference: "https://datatracker.ietf.org/doc/html/rfc7845#section-5.2.1"}
			}
		}
	default:
		return Representation{Class: Unsupported}
	}
	if key == "" {
		if format == MP3v1 || !d.Mutable {
			return Representation{Class: Unsupported, Reference: ref}
		}
		return Representation{Class: Fallback, Key: fallbackKey(format, d.ID), Note: "Lossless AudioMeta encoding; external tools may expose the raw value."}
	}
	switch key[0] {
	case '~':
		class, key = Conventional, key[1:]
	case '?':
		class, key = Partial, key[1:]
	case '!':
		class, key = ReadOnly, key[1:]
	}
	names := strings.Split(key, ",")
	result := Representation{Class: class, Key: names[0], Aliases: names[1:], Reference: ref}
	if class == Conventional && strings.HasPrefix(names[0], "TXXX:") {
		result.Reference = "https://picard-docs.musicbrainz.org/en/latest/appendices/tag_mapping.html"
	}
	if format == FLAC || format == Vorbis || format == Opus {
		switch d.ID {
		case RecordingDate:
			result.Note = "Xiph defines DATE as recording date; Picard uses DATE as a release-date convention. AudioMeta reads DATE as recording date and uses RELEASEDATE separately."
			result.Reference = "https://xiph.org/vorbis/doc/v-comment.html"
		case Rating:
			result.Note = "RATING uses the normalized 0–1 scale. Higher precision and user-qualified keys stay opaque."
			result.Reference = "https://picard-docs.musicbrainz.org/en/latest/config/options_ratings.html"
		case Title, Album, Artist, Performer, Copyright, License, Organization, Description, Genre, Location, Contact, ISRC, TrackNumber, Version:
			result.Reference = "https://xiph.org/vorbis/doc/v-comment.html"
		}
	}
	if class == Partial {
		result.Note = "Native representation is version-, shape-, or context-dependent; preserve opaque values when projection is incomplete."
	}
	return result
}

// FallbackMapping returns the lossless private representation for a mutable
// field and supported writable format.
func (d Definition) FallbackMapping(format Format) Representation {
	if !d.Mutable {
		return Representation{Class: Unsupported}
	}
	switch format {
	case MP3v22, MP3v23, MP3v24, MP4, FLAC, Vorbis, Opus:
		return Representation{Class: Fallback, Key: fallbackKey(format, d.ID)}
	default:
		return Representation{Class: Unsupported}
	}
}

func fallbackKey(format Format, field Field) string {
	switch format {
	case MP3v22:
		return "TXX:" + FallbackName(field)
	case MP3v23, MP3v24:
		return "TXXX:" + FallbackName(field)
	case MP4:
		return "----:" + MP4Namespace + ":" + FallbackName(field)
	default:
		return CommentFallbackName(field)
	}
}

const MP4Namespace = "io.github.gcottom.audiometa"

func FallbackName(field Field) string { return "audiometa/1/" + string(field) }
func CustomName(key string) string    { return "audiometa/1/custom/" + key }
func ParseCustomName(name string) (string, bool) {
	if !strings.HasPrefix(name, "audiometa/1/custom/") {
		return "", false
	}
	key, err := NormalizeCustomKey(strings.TrimPrefix(name, "audiometa/1/custom/"))
	return key, err == nil
}
func CommentFallbackName(field Field) string { return "AUDIOMETA_1_" + strings.ToUpper(string(field)) }

const commentCustomPrefix = "AUDIOMETA_1_CUSTOM_"

func Fields() []Definition                  { return append([]Definition(nil), definitions...) }
func Lookup(field Field) (Definition, bool) { definition, ok := registry[field]; return definition, ok }

func definition(id Field, name string, kind Kind, multiple bool, v, i, m, i3, i2 string) Definition {
	d := Definition{ID: id, Name: name, Meaning: name, Kind: kind, Multiple: multiple, Mutable: true, comment: v, id3v24: i, mp4: m, id3v23: i3, id3v22: i2}
	d.Constraints = kindConstraints[kind]
	switch id {
	case ISRC:
		d.Constraints = "ISO 3901 identifier in 12-character compact form."
	case UPC:
		d.Constraints = "12 decimal digits."
	case Barcode:
		d.Constraints = "8, 12, 13 or 14 decimal digits."
	case ASIN:
		d.Constraints = "10 uppercase letters or digits."
	case AudioDuration, PlaylistDelay:
		d.Unit = "ms"
	case EncoderDelay, EncoderPadding:
		d.Unit = "samples"
	case BPM:
		d.Unit = "beats/minute"
	case Rating:
		d.Unit, d.Meaning, d.Constraints = "percent", "User rating on a 0–100 scale", "Integer from 0 through 100 inclusive."
	case ReplayGainTrackGain, ReplayGainAlbumGain:
		d.Unit = "dB"
	case ReplayGainReferenceLoudness:
		d.Unit = "dB SPL"
	case ReplayGainTrackPeak, ReplayGainAlbumPeak:
		d.Unit, d.Constraints = "linear amplitude ratio", "Nonnegative finite decimal; may exceed 1."
	case R128TrackGain, R128AlbumGain:
		d.Unit, d.Constraints = "signed Q7.8 dB", "Integer from −32768 through 32767 inclusive, relative to OpusHead output gain."
	case IntegratedLoudness:
		d.Unit = "LUFS"
	case LoudnessRange:
		d.Unit, d.Constraints = "LU", "Nonnegative finite decimal."
	case Artist:
		d.Meaning = "Track artist credits in source order."
	case Artists:
		d.Meaning = "Explicit ordered individual artists, independent of display credit."
	case ArtistDisplay:
		d.Meaning = "Literal display credit for the track artist."
	case RecordingID:
		d.Meaning = "MusicBrainz recording identifier; historically called Track Id."
	case ReleaseTrackID:
		d.Meaning = "MusicBrainz track identifier within a release."
	case RecordingDate:
		d.Meaning = "When the audio was recorded."
	case ReleaseDate:
		d.Meaning = "Release date of the specific edition."
	case OriginalReleaseDate:
		d.Meaning = "Original release date of the work or recording."
	case GaplessAlbum:
		d.Meaning = "Album-level gapless playback preference."
	}
	switch id {
	case AudioDuration, EncoderDelay, EncoderPadding:
		d.Mutable = false
	}
	return d
}

var definitions = []Definition{
	definition(Title, "Title", KindText, false, "TITLE", "TIT2", "©nam", "TIT2", "TT2"),
	definition(Subtitle, "Subtitle", KindText, false, "SUBTITLE", "TIT3", "~----:com.apple.iTunes:SUBTITLE", "TIT3", "TT3"),
	definition(Album, "Album", KindText, false, "ALBUM", "TALB", "©alb", "TALB", "TAL"),
	definition(Grouping, "Grouping", KindText, false, "GROUPING", "TIT1", "©grp", "TIT1", "TT1"),
	definition(Version, "Version", KindText, false, "VERSION", "~TXXX:VERSION", "", "", ""),
	definition(OriginalTitle, "Original Title", KindText, false, "^ORIGINALTITLE", "", "", "", ""),
	definition(Artist, "Artist", KindText, true, "ARTIST", "TPE1", "©ART", "?TPE1", "?TP1"),
	definition(ArtistDisplay, "Artist Display", KindText, false, "ARTIST_CREDIT", "~TXXX:ARTIST_CREDIT", "", "", ""),
	definition(Artists, "Artists", KindText, true, "ARTISTS", "~TXXX:ARTISTS", "~----:com.apple.iTunes:ARTISTS", "", ""),
	definition(AlbumArtist, "Album Artist", KindText, true, "ALBUMARTIST,ALBUM ARTIST", "~TPE2", "aART", "?TPE2", "?TP2"),
	definition(AlbumArtistDisplay, "Album Artist Display", KindText, false, "ALBUMARTIST_CREDIT", "~TXXX:ALBUMARTIST_CREDIT", "", "", ""),
	definition(OriginalArtist, "Original Artist", KindText, true, "ORIGINALARTIST", "TOPE", "~----:com.apple.iTunes:ORIGINALARTIST", "?TOPE", "?TOA"),
	definition(Performer, "Performer", KindText, true, "PERFORMER", "~TXXX:PERFORMER", "~----:com.apple.iTunes:PERFORMER", "", ""),
	definition(FeaturedArtist, "Featured Artist", KindText, true, "FEATUREDARTIST", "~TXXX:FEATUREDARTIST", "", "", ""),
	definition(PerformerCredits, "Performer Credits", KindCredit, true, "", "?TMCL", "", "", ""),
	definition(Composer, "Composer", KindText, true, "COMPOSER", "TCOM", "©wrt", "?TCOM", "?TCM"),
	definition(Lyricist, "Lyricist", KindText, true, "LYRICIST", "TEXT", "~----:com.apple.iTunes:LYRICIST", "?TEXT", "?TXT"),
	definition(Conductor, "Conductor", KindText, true, "CONDUCTOR", "TPE3", "~----:com.apple.iTunes:CONDUCTOR", "?TPE3", "?TP3"),
	definition(Arranger, "Arranger", KindText, true, "ARRANGER", "~TXXX:ARRANGER", "~----:com.apple.iTunes:ARRANGER", "", ""),
	definition(Work, "Work", KindText, false, "WORK", "~TXXX:WORK", "©wrk", "~TXXX:WORK", ""),
	definition(Movement, "Movement", KindText, false, "MOVEMENTNAME", "~MVNM", "©mvn", "~MVNM", ""),
	definition(MovementNumber, "Movement Number", KindInt, false, "MOVEMENT,MOVEMENTNUMBER", "~MVIN", "©mvi", "~MVIN", ""),
	definition(MovementTotal, "Movement Total", KindInt, false, "MOVEMENTTOTAL", "~MVIN", "©mvc", "~MVIN", ""),
	definition(ShowMovement, "Show Movement", KindBool, false, "SHOWMOVEMENT", "~TXXX:SHOWMOVEMENT", "shwm", "", ""),
	definition(OpusNumber, "Opus Number", KindText, false, "OPUS", "~TXXX:OPUS", "", "", ""),
	definition(Producer, "Producer", KindText, true, "PRODUCER", "~TXXX:PRODUCER", "~----:com.apple.iTunes:PRODUCER", "", ""),
	definition(Engineer, "Engineer", KindText, true, "ENGINEER", "~TXXX:ENGINEER", "~----:com.apple.iTunes:ENGINEER", "", ""),
	definition(Mixer, "Mixer", KindText, true, "MIXER", "~TXXX:MIXER", "~----:com.apple.iTunes:MIXER", "", ""),
	definition(Remixer, "Remixer", KindText, true, "REMIXER", "TPE4", "~----:com.apple.iTunes:REMIXER", "?TPE4", "?TP4"),
	definition(MasteringEngineer, "Mastering Engineer", KindText, true, "MASTERINGENGINEER", "~TXXX:MASTERINGENGINEER", "~----:com.apple.iTunes:MASTERINGENGINEER", "", ""),
	definition(ProductionCredits, "Production Credits", KindCredit, true, "", "?TIPL", "", "?IPLS", ""),
	definition(Publisher, "Publisher", KindText, false, "PUBLISHER", "TPUB", "~----:com.apple.iTunes:PUBLISHER", "TPUB", "TPB"),
	definition(RecordLabel, "Record Label", KindText, true, "LABEL", "~TXXX:LABEL", "~----:com.apple.iTunes:LABEL", "", ""),
	definition(Organization, "Organization", KindText, false, "ORGANIZATION", "~TXXX:ORGANIZATION", "", "", ""),
	definition(CatalogNumber, "Catalog Number", KindText, true, "CATALOGNUMBER", "~TXXX:CATALOGNUMBER", "~----:com.apple.iTunes:CATALOGNUMBER", "", ""),
	definition(TrackNumber, "Track Number", KindInt, false, "TRACKNUMBER", "TRCK", "trkn", "TRCK", "TRK"),
	definition(TrackTotal, "Track Total", KindInt, false, "TRACKTOTAL,TOTALTRACKS", "TRCK", "trkn", "TRCK", "TRK"),
	definition(DiscNumber, "Disc Number", KindInt, false, "DISCNUMBER", "TPOS", "disk", "TPOS", "TPA"),
	definition(DiscTotal, "Disc Total", KindInt, false, "DISCTOTAL,TOTALDISCS", "TPOS", "disk", "TPOS", "TPA"),
	definition(Compilation, "Compilation", KindBool, false, "COMPILATION", "~TCMP", "cpil", "~TCMP", "~TCP"),
	definition(MediaType, "Media Type", KindText, false, "MEDIA,MEDIATYPE", "TMED", "~----:com.apple.iTunes:MEDIA", "TMED", "TMT"),
	definition(MediaKind, "Media Kind", KindInt, false, "MEDIAKIND", "~TXXX:MEDIAKIND", "stik", "", ""),
	definition(Genre, "Genre", KindText, true, "GENRE", "TCON", "©gen,gnre", "?TCON", "?TCO"),
	definition(Mood, "Mood", KindText, true, "MOOD", "TMOO", "~----:com.apple.iTunes:MOOD", "~TXXX:MOOD", ""),
	definition(BPM, "BPM", KindInt, false, "BPM", "TBPM", "tmpo", "TBPM", "TBP"),
	definition(MusicalKey, "Musical Key", KindText, false, "INITIALKEY,KEY", "TKEY", "~----:com.apple.iTunes:initialkey", "TKEY", "TKE"),
	definition(Language, "Language", KindText, true, "LANGUAGE", "TLAN", "~----:com.apple.iTunes:LANGUAGE", "?TLAN", "?TLA"),
	definition(RecordingDate, "Recording Date", KindDate, false, "DATE,RECORDINGDATE", "TDRC", "~----:com.apple.iTunes:RECORDINGDATE", "?TYER", "?TYE"),
	definition(ReleaseDate, "Release Date", KindDate, false, "RELEASEDATE", "TDRL", "©day", "~TXXX:RELEASEDATE", ""),
	definition(OriginalReleaseDate, "Original Release Date", KindDate, false, "ORIGINALDATE,ORIGINALYEAR", "TDOR", "~----:com.apple.iTunes:ORIGINALDATE", "?TORY", "?TOR"),
	definition(EncodingDate, "Encoding Date", KindDate, false, "ENCODINGDATE", "TDEN", "", "", ""),
	definition(TaggingDate, "Tagging Date", KindDate, false, "TAGGINGDATE", "TDTG", "", "", ""),
	definition(Copyright, "Copyright", KindText, false, "COPYRIGHT", "TCOP", "cprt,©cpy", "TCOP", "TCR"),
	definition(License, "License", KindText, false, "LICENSE", "~TXXX:LICENSE", "~----:com.apple.iTunes:LICENSE", "", ""),
	definition(Ownership, "Ownership", KindText, false, "^OWNERSHIP", "?OWNE", "", "?OWNE", ""),
	definition(PaymentURL, "Payment URL", KindURL, false, "PAYMENTURL", "WPAY", "", "WPAY", ""),
	definition(TermsOfUse, "Terms Of Use", KindText, false, "^TERMSOFUSE", "?USER", "", "?USER", ""),
	definition(ISRC, "ISRC", KindText, true, "ISRC", "TSRC", "~----:com.apple.iTunes:ISRC", "?TSRC", "?TRC"),
	definition(UPC, "UPC", KindText, false, "UPC", "~TXXX:UPC", "~----:com.apple.iTunes:UPC", "", ""),
	definition(Barcode, "Barcode", KindText, false, "BARCODE", "~TXXX:BARCODE", "~----:com.apple.iTunes:BARCODE", "", ""),
	definition(ASIN, "ASIN", KindText, false, "ASIN", "~TXXX:ASIN", "~----:com.apple.iTunes:ASIN", "", ""),
	definition(RecordingID, "MusicBrainz Recording ID", KindUUID, false, "MUSICBRAINZ_TRACKID", "~UFID:http://musicbrainz.org", "~----:com.apple.iTunes:MusicBrainz Track Id", "~UFID:http://musicbrainz.org", "~UFI:http://musicbrainz.org"),
	definition(ReleaseTrackID, "MusicBrainz Release Track ID", KindUUID, false, "MUSICBRAINZ_RELEASETRACKID", "~TXXX:MusicBrainz Release Track Id", "~----:com.apple.iTunes:MusicBrainz Release Track Id", "", ""),
	definition(ReleaseID, "MusicBrainz Release ID", KindUUID, false, "MUSICBRAINZ_ALBUMID", "~TXXX:MusicBrainz Album Id", "~----:com.apple.iTunes:MusicBrainz Album Id", "", ""),
	definition(ReleaseGroupID, "MusicBrainz Release Group ID", KindUUID, false, "MUSICBRAINZ_RELEASEGROUPID", "~TXXX:MusicBrainz Release Group Id", "~----:com.apple.iTunes:MusicBrainz Release Group Id", "", ""),
	definition(ArtistID, "MusicBrainz Artist ID", KindUUID, true, "MUSICBRAINZ_ARTISTID", "~TXXX:MusicBrainz Artist Id", "~----:com.apple.iTunes:MusicBrainz Artist Id", "", ""),
	definition(AlbumArtistID, "MusicBrainz Album Artist ID", KindUUID, true, "MUSICBRAINZ_ALBUMARTISTID", "~TXXX:MusicBrainz Album Artist Id", "~----:com.apple.iTunes:MusicBrainz Album Artist Id", "", ""),
	definition(WorkID, "MusicBrainz Work ID", KindUUID, true, "MUSICBRAINZ_WORKID", "~TXXX:MusicBrainz Work Id", "~----:com.apple.iTunes:MusicBrainz Work Id", "", ""),
	definition(DiscID, "MusicBrainz Disc ID", KindText, false, "MUSICBRAINZ_DISCID", "~TXXX:MusicBrainz Disc Id", "~----:com.apple.iTunes:MusicBrainz Disc Id", "", ""),
	definition(AcoustID, "AcoustID ID", KindUUID, false, "ACOUSTID_ID", "~TXXX:Acoustid Id", "~----:com.apple.iTunes:Acoustid Id", "", ""),
	definition(AcoustIDFingerprint, "AcoustID Fingerprint", KindText, false, "ACOUSTID_FINGERPRINT", "~TXXX:Acoustid Fingerprint", "~----:com.apple.iTunes:Acoustid Fingerprint", "", ""),
	definition(MusicIPFingerprint, "Musicip Fingerprint", KindText, false, "MUSICIP_PUID", "~TXXX:MusicIP PUID", "~----:com.apple.iTunes:MusicIP PUID", "", ""),
	definition(BinaryIdentifier, "Binary Identifier", KindBinary, false, "", "", "", "", ""),
	definition(TitleSort, "Title Sort", KindText, false, "TITLESORT", "TSOT", "sonm", "~TSOT", "~TST"),
	definition(ArtistSort, "Artist Sort", KindText, false, "ARTISTSORT", "TSOP", "soar", "~TSOP", "~TSP"),
	definition(AlbumSort, "Album Sort", KindText, false, "ALBUMSORT", "TSOA", "soal", "~TSOA", "~TSA"),
	definition(AlbumArtistSort, "Album Artist Sort", KindText, false, "ALBUMARTISTSORT", "~TSO2", "soaa", "~TSO2", "~TS2"),
	definition(ComposerSort, "Composer Sort", KindText, false, "COMPOSERSORT", "~TSOC", "soco", "~TSOC", "~TSC"),
	definition(Comments, "Comments", KindText, true, "COMMENT", "?COMM", "©cmt", "?COMM", "?COM"),
	definition(Description, "Description", KindText, false, "DESCRIPTION", "~TXXX:DESCRIPTION", "desc", "", ""),
	definition(LongDescription, "Long Description", KindText, false, "LONGDESCRIPTION", "~TXXX:LONGDESCRIPTION", "ldes", "", ""),
	definition(Lyrics, "Lyrics", KindText, true, "LYRICS,UNSYNCEDLYRICS", "?USLT", "©lyr", "?USLT", "?ULT"),
	definition(SynchronizedLyrics, "Synchronized Lyrics", KindSyncedLyrics, true, "", "?SYLT", "", "?SYLT", "?SLT"),
	definition(Chapters, "Chapters", KindChapter, true, "", "?CHAP", "?chpl", "?CHAP", ""),
	definition(TimedMetadata, "Timed Metadata", KindTimedEvent, true, "", "", "", "", ""),
	definition(Podcast, "Podcast", KindBool, false, "^PODCAST", "~PCST", "pcst", "~PCST", ""),
	definition(PodcastShow, "Podcast Show", KindText, false, "^PODCASTSHOW", "~TXXX:PODCASTSHOW", "tvsh", "", ""),
	definition(PodcastEpisode, "Podcast Episode", KindInt, false, "^PODCASTEPISODE", "~TXXX:PODCASTEPISODE", "tves", "", ""),
	definition(PodcastSeason, "Podcast Season", KindInt, false, "^PODCASTSEASON", "~TXXX:PODCASTSEASON", "tvsn", "", ""),
	definition(PodcastCategory, "Podcast Category", KindText, false, "^PODCASTCATEGORY", "~TCAT", "catg", "~TCAT", ""),
	definition(PodcastFeedURL, "Podcast Feed URL", KindURL, false, "^PODCASTURL", "~WFED", "purl", "~WFED", ""),
	definition(PodcastEpisodeID, "Podcast Episode ID", KindText, false, "^PODCASTID", "~TGID", "egid", "~TGID", ""),
	definition(Artwork, "Artwork", KindPicture, true, "METADATA_BLOCK_PICTURE", "APIC", "?covr", "APIC", "PIC"),
	definition(ReplayGainTrackGain, "Replaygain Track Gain", KindDecimal, false, "REPLAYGAIN_TRACK_GAIN", "~TXXX:REPLAYGAIN_TRACK_GAIN", "~----:com.apple.iTunes:REPLAYGAIN_TRACK_GAIN,----:com.apple.iTunes:replaygain_track_gain", "", ""),
	definition(ReplayGainAlbumGain, "Replaygain Album Gain", KindDecimal, false, "REPLAYGAIN_ALBUM_GAIN", "~TXXX:REPLAYGAIN_ALBUM_GAIN", "~----:com.apple.iTunes:REPLAYGAIN_ALBUM_GAIN,----:com.apple.iTunes:replaygain_album_gain", "", ""),
	definition(ReplayGainTrackPeak, "Replaygain Track Peak", KindDecimal, false, "REPLAYGAIN_TRACK_PEAK", "~TXXX:REPLAYGAIN_TRACK_PEAK", "~----:com.apple.iTunes:REPLAYGAIN_TRACK_PEAK,----:com.apple.iTunes:replaygain_track_peak", "", ""),
	definition(ReplayGainAlbumPeak, "Replaygain Album Peak", KindDecimal, false, "REPLAYGAIN_ALBUM_PEAK", "~TXXX:REPLAYGAIN_ALBUM_PEAK", "~----:com.apple.iTunes:REPLAYGAIN_ALBUM_PEAK,----:com.apple.iTunes:replaygain_album_peak", "", ""),
	definition(ReplayGainReferenceLoudness, "Replaygain Reference Loudness", KindDecimal, false, "REPLAYGAIN_REFERENCE_LOUDNESS", "~TXXX:REPLAYGAIN_REFERENCE_LOUDNESS", "~----:com.apple.iTunes:REPLAYGAIN_REFERENCE_LOUDNESS,----:com.apple.iTunes:replaygain_reference_loudness", "", ""),
	definition(R128TrackGain, "R128 Track Gain", KindInt, false, "", "", "", "", ""),
	definition(R128AlbumGain, "R128 Album Gain", KindInt, false, "", "", "", "", ""),
	definition(IntegratedLoudness, "Integrated Loudness", KindDecimal, false, "", "", "", "", ""),
	definition(LoudnessRange, "Loudness Range", KindDecimal, false, "", "", "", "", ""),
	definition(Rating, "Rating", KindInt, false, "RATING", "?POPM", "", "?POPM", "?POP"),
	definition(ContentRating, "Content Rating", KindInt, false, "CONTENTRATING", "~TXXX:CONTENTRATING", "rtng", "", ""),
	definition(PlayCount, "Play Count", KindInt, false, "PLAYCOUNT", "PCNT", "", "PCNT", "CNT"),
	definition(GaplessAlbum, "Gapless Album", KindBool, false, "GAPLESS", "~TXXX:GAPLESS", "pgap", "", ""),
	definition(OfficialArtistURL, "Official Artist URL", KindURL, true, "WEBSITE", "WOAR", "~----:com.apple.iTunes:WEBSITE", "WOAR", "WAR"),
	definition(PublisherURL, "Publisher URL", KindURL, false, "PUBLISHERURL", "WPUB", "", "WPUB", "WPB"),
	definition(PurchaseURL, "Purchase URL", KindURL, false, "PURCHASEURL", "WCOM", "", "WCOM", "WCM"),
	definition(CopyrightURL, "Copyright URL", KindURL, false, "COPYRIGHTURL", "WCOP", "", "WCOP", "WCP"),
	definition(AudioSourceURL, "Audio Source URL", KindURL, false, "AUDIOSOURCEURL", "WOAS", "", "WOAS", "WAS"),
	definition(OfficialAudioURL, "Official Audio URL", KindURL, false, "OFFICIALAUDIOURL", "WOAF", "", "WOAF", "WAF"),
	definition(Encoder, "Encoder", KindText, false, "ENCODER", "TENC", "©too", "TENC", "TEN"),
	definition(EncodingSettings, "Encoding Settings", KindText, false, "ENCODERSETTINGS", "TSSE", "~----:com.apple.iTunes:ENCODERSETTINGS", "TSSE", "TSS"),
	definition(MediaOwner, "Media Owner", KindText, false, "^MEDIAOWNER", "TOWN", "", "TOWN", ""),
	definition(Contact, "Contact", KindText, false, "CONTACT", "~TXXX:CONTACT", "~----:com.apple.iTunes:CONTACT", "", ""),
	definition(Location, "Location", KindText, false, "LOCATION", "~TXXX:LOCATION", "~----:com.apple.iTunes:LOCATION", "", ""),
	definition(ReleaseStatus, "Release Status", KindText, false, "RELEASESTATUS", "~TXXX:MusicBrainz Album Status", "~----:com.apple.iTunes:MusicBrainz Album Status", "", ""),
	definition(ReleaseType, "Release Type", KindText, true, "RELEASETYPE", "~TXXX:MusicBrainz Album Type", "~----:com.apple.iTunes:MusicBrainz Album Type", "", ""),
	definition(ReleaseCountry, "Release Country", KindText, false, "RELEASECOUNTRY", "~TXXX:MusicBrainz Album Release Country", "~----:com.apple.iTunes:MusicBrainz Album Release Country", "", ""),
	definition(Script, "Script", KindText, false, "SCRIPT", "~TXXX:SCRIPT", "~----:com.apple.iTunes:SCRIPT", "", ""),
	definition(DiscSubtitle, "Disc Subtitle", KindText, false, "DISCSUBTITLE", "TSST", "~----:com.apple.iTunes:DISCSUBTITLE", "~TXXX:DISCSUBTITLE", ""),
	definition(OriginalFilename, "Original Filename", KindText, false, "ORIGINALFILENAME", "TOFN", "", "TOFN", "TOF"),
	definition(OriginalLyricist, "Original Lyricist", KindText, true, "ORIGINALLYRICIST", "TOLY", "", "?TOLY", "?TOL"),
	definition(PlaylistDelay, "Playlist Delay", KindInt, false, "PLAYLISTDELAY", "TDLY", "", "TDLY", "TDY"),
	definition(RadioStationName, "Radio Station Name", KindText, false, "RADIOSTATION", "TRSN", "", "TRSN", ""),
	definition(RadioStationOwner, "Radio Station Owner", KindText, false, "RADIOSTATIONOWNER", "TRSO", "", "TRSO", ""),
	definition(ISWC, "ISWC", KindText, true, "ISWC", "~TXXX:ISWC", "~----:com.apple.iTunes:ISWC", "", ""),
	definition(OriginalAlbum, "Original Album", KindText, false, "ORIGINALALBUM", "TOAL", "~----:com.apple.iTunes:ORIGINALALBUM", "TOAL", "TOT"),
	definition(InitialProducer, "Initial Producer", KindText, false, "INITIALPRODUCER", "TPRO", "", "", ""),
	definition(AudioDuration, "Audio Duration", KindInt, false, "", "!TLEN", "", "!TLEN", "!TLE"),
	definition(EncoderDelay, "Encoder Delay", KindInt, false, "", "", "!----:com.apple.iTunes:iTunSMPB", "", ""),
	definition(EncoderPadding, "Encoder Padding", KindInt, false, "", "", "!----:com.apple.iTunes:iTunSMPB", "", ""),
}

var registry = func() map[Field]Definition {
	byID := make(map[Field]Definition, len(definitions))
	for _, d := range definitions {
		byID[d.ID] = d
	}
	return byID
}()
