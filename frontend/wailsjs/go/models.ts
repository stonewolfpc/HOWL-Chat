export namespace audio {
	
	export class DesiredSettings {
	    asrEnabled: boolean;
	    ttsEnabled: boolean;
	    asrEngine: string;
	    asrLanguage: string;
	    asrDevice: string;
	    asrThreads: number;
	    asrBeamSize: number;
	    asrBestOf: number;
	    asrTemperature: number;
	    asrTranslateToEnglish: boolean;
	    asrNoTimestamps: boolean;
	    asrFlashAttention: boolean;
	    vadEnabled: boolean;
	    vadThreshold: number;
	    vadMinSpeechMs: number;
	    vadMinSilenceMs: number;
	    vadPaddingMs: number;
	    ttsEngine: string;
	    ttsVoice: string;
	    ttsSpeed: number;
	    ttsVolume: number;
	    ttsNoiseScale: number;
	    ttsNoiseW: number;
	    ttsSentenceSilence: number;
	    ttsAutoplayMessage: boolean;
	    ttsCacheEnabled: boolean;
	    ttsCacheSizeMB: number;
	    audioInputDevice: string;
	    audioOutputDevice: string;
	    audioSampleRate: number;
	    audioChannels: string;
	    audioBitDepth: number;
	    micSensitivity: number;
	    pushToTalk: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DesiredSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.asrEnabled = source["asrEnabled"];
	        this.ttsEnabled = source["ttsEnabled"];
	        this.asrEngine = source["asrEngine"];
	        this.asrLanguage = source["asrLanguage"];
	        this.asrDevice = source["asrDevice"];
	        this.asrThreads = source["asrThreads"];
	        this.asrBeamSize = source["asrBeamSize"];
	        this.asrBestOf = source["asrBestOf"];
	        this.asrTemperature = source["asrTemperature"];
	        this.asrTranslateToEnglish = source["asrTranslateToEnglish"];
	        this.asrNoTimestamps = source["asrNoTimestamps"];
	        this.asrFlashAttention = source["asrFlashAttention"];
	        this.vadEnabled = source["vadEnabled"];
	        this.vadThreshold = source["vadThreshold"];
	        this.vadMinSpeechMs = source["vadMinSpeechMs"];
	        this.vadMinSilenceMs = source["vadMinSilenceMs"];
	        this.vadPaddingMs = source["vadPaddingMs"];
	        this.ttsEngine = source["ttsEngine"];
	        this.ttsVoice = source["ttsVoice"];
	        this.ttsSpeed = source["ttsSpeed"];
	        this.ttsVolume = source["ttsVolume"];
	        this.ttsNoiseScale = source["ttsNoiseScale"];
	        this.ttsNoiseW = source["ttsNoiseW"];
	        this.ttsSentenceSilence = source["ttsSentenceSilence"];
	        this.ttsAutoplayMessage = source["ttsAutoplayMessage"];
	        this.ttsCacheEnabled = source["ttsCacheEnabled"];
	        this.ttsCacheSizeMB = source["ttsCacheSizeMB"];
	        this.audioInputDevice = source["audioInputDevice"];
	        this.audioOutputDevice = source["audioOutputDevice"];
	        this.audioSampleRate = source["audioSampleRate"];
	        this.audioChannels = source["audioChannels"];
	        this.audioBitDepth = source["audioBitDepth"];
	        this.micSensitivity = source["micSensitivity"];
	        this.pushToTalk = source["pushToTalk"];
	    }
	}

}

export namespace chat {
	
	export class Image {
	    path: string;
	    thumbnail?: string;
	
	    static createFrom(source: any = {}) {
	        return new Image(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.thumbnail = source["thumbnail"];
	    }
	}
	export class CharacterContext {
	    id: string;
	    name: string;
	    description: string;
	    personality: string;
	    backstory: string;
	    goals: string[];
	    relationships: string[];
	    references: string[];
	    age?: string;
	    role?: string;
	    background?: string;
	    priority?: number;
	    version?: number;
	    voice_config?: Record<string, any>;
	    image?: Image;
	
	    static createFrom(source: any = {}) {
	        return new CharacterContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.personality = source["personality"];
	        this.backstory = source["backstory"];
	        this.goals = source["goals"];
	        this.relationships = source["relationships"];
	        this.references = source["references"];
	        this.age = source["age"];
	        this.role = source["role"];
	        this.background = source["background"];
	        this.priority = source["priority"];
	        this.version = source["version"];
	        this.voice_config = source["voice_config"];
	        this.image = this.convertValues(source["image"], Image);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ScenarioContextCharacter {
	    id: string;
	    name: string;
	    role?: string;
	    avatar?: string;
	
	    static createFrom(source: any = {}) {
	        return new ScenarioContextCharacter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.role = source["role"];
	        this.avatar = source["avatar"];
	    }
	}
	export class ScenarioContext {
	    id: string;
	    name: string;
	    description: string;
	    key_details: string[];
	    references: string[];
	    image?: Image;
	    characters?: ScenarioContextCharacter[];
	    lorebooks?: any[];
	
	    static createFrom(source: any = {}) {
	        return new ScenarioContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.key_details = source["key_details"];
	        this.references = source["references"];
	        this.image = this.convertValues(source["image"], Image);
	        this.characters = this.convertValues(source["characters"], ScenarioContextCharacter);
	        this.lorebooks = source["lorebooks"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WorldContext {
	    id: string;
	    name: string;
	    description: string;
	    key_facts: string[];
	    references: string[];
	    image?: Image;
	
	    static createFrom(source: any = {}) {
	        return new WorldContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.key_facts = source["key_facts"];
	        this.references = source["references"];
	        this.image = this.convertValues(source["image"], Image);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChatHistory {
	    id: string;
	    sender: string;
	    text: string;
	    role: string;
	    timestamp: string;
	    hasImage: boolean;
	    characterId?: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sender = source["sender"];
	        this.text = source["text"];
	        this.role = source["role"];
	        this.timestamp = source["timestamp"];
	        this.hasImage = source["hasImage"];
	        this.characterId = source["characterId"];
	    }
	}
	export class ChatMeta {
	    id: string;
	    name: string;
	    createdAt: string;
	    updatedAt: string;
	    worldName: string;
	    scenarioName: string;
	    characterNames: string[];
	    lorebookNames: string[];
	    worldId: string;
	    scenarioId: string;
	    activeCharacters: string[];
	    removedCharacters: string[];
	    addedCharacters: string[];
	    activeLorebooks: string[];
	    removedLorebooks: string[];
	    addedLorebooks: string[];
	
	    static createFrom(source: any = {}) {
	        return new ChatMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.worldName = source["worldName"];
	        this.scenarioName = source["scenarioName"];
	        this.characterNames = source["characterNames"];
	        this.lorebookNames = source["lorebookNames"];
	        this.worldId = source["worldId"];
	        this.scenarioId = source["scenarioId"];
	        this.activeCharacters = source["activeCharacters"];
	        this.removedCharacters = source["removedCharacters"];
	        this.addedCharacters = source["addedCharacters"];
	        this.activeLorebooks = source["activeLorebooks"];
	        this.removedLorebooks = source["removedLorebooks"];
	        this.addedLorebooks = source["addedLorebooks"];
	    }
	}
	export class ChatData {
	    meta: ChatMeta;
	    history: ChatHistory[];
	    world?: WorldContext;
	    scenario?: ScenarioContext;
	    characters: CharacterContext[];
	    lorebooks: any[];
	    memories: any[];
	
	    static createFrom(source: any = {}) {
	        return new ChatData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], ChatMeta);
	        this.history = this.convertValues(source["history"], ChatHistory);
	        this.world = this.convertValues(source["world"], WorldContext);
	        this.scenario = this.convertValues(source["scenario"], ScenarioContext);
	        this.characters = this.convertValues(source["characters"], CharacterContext);
	        this.lorebooks = source["lorebooks"];
	        this.memories = source["memories"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class ChatSummary {
	    id: string;
	    name: string;
	    createdAt: string;
	    updatedAt: string;
	    worldName: string;
	    scenarioName: string;
	    characterNames: string[];
	    lorebookNames: string[];
	    messageCount: number;
	
	    static createFrom(source: any = {}) {
	        return new ChatSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.worldName = source["worldName"];
	        this.scenarioName = source["scenarioName"];
	        this.characterNames = source["characterNames"];
	        this.lorebookNames = source["lorebookNames"];
	        this.messageCount = source["messageCount"];
	    }
	}
	
	
	

}

export namespace contentimport {
	
	export class ItemResult {
	    source_key: string;
	    item_id?: string;
	    state: string;
	    reason?: string;
	    path?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ItemResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_key = source["source_key"];
	        this.item_id = source["item_id"];
	        this.state = source["state"];
	        this.reason = source["reason"];
	        this.path = source["path"];
	        this.message = source["message"];
	    }
	}
	export class JournalItem {
	    source_key: string;
	    source_kind: string;
	    source_id?: string;
	    source_name?: string;
	    source_root?: string;
	    folder?: string;
	    item_id?: string;
	    library_id?: string;
	    root?: string;
	    final_path?: string;
	    staging_path?: string;
	    owner_path?: string;
	    revision?: string;
	    plan_hash?: string;
	    state: string;
	    reason?: string;
	    message?: string;
	    evidence?: string;
	    source_hash?: string;
	    document_hash?: string;
	    lore_hashes?: Record<string, string>;
	    media_hash?: string;
	
	    static createFrom(source: any = {}) {
	        return new JournalItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_key = source["source_key"];
	        this.source_kind = source["source_kind"];
	        this.source_id = source["source_id"];
	        this.source_name = source["source_name"];
	        this.source_root = source["source_root"];
	        this.folder = source["folder"];
	        this.item_id = source["item_id"];
	        this.library_id = source["library_id"];
	        this.root = source["root"];
	        this.final_path = source["final_path"];
	        this.staging_path = source["staging_path"];
	        this.owner_path = source["owner_path"];
	        this.revision = source["revision"];
	        this.plan_hash = source["plan_hash"];
	        this.state = source["state"];
	        this.reason = source["reason"];
	        this.message = source["message"];
	        this.evidence = source["evidence"];
	        this.source_hash = source["source_hash"];
	        this.document_hash = source["document_hash"];
	        this.lore_hashes = source["lore_hashes"];
	        this.media_hash = source["media_hash"];
	    }
	}
	export class Operation {
	    version: number;
	    category: string;
	    operation_id: string;
	    plan_hash: string;
	    destination_library_id: string;
	    destination_root?: string;
	    status: string;
	    started_at: string;
	    updated_at: string;
	    items: JournalItem[];
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new Operation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.category = source["category"];
	        this.operation_id = source["operation_id"];
	        this.plan_hash = source["plan_hash"];
	        this.destination_library_id = source["destination_library_id"];
	        this.destination_root = source["destination_root"];
	        this.status = source["status"];
	        this.started_at = source["started_at"];
	        this.updated_at = source["updated_at"];
	        this.items = this.convertValues(source["items"], JournalItem);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ReferenceIssue {
	    field: string;
	    index: number;
	    reason: string;
	    kind?: string;
	    id?: string;
	    legacy_name?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ReferenceIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.field = source["field"];
	        this.index = source["index"];
	        this.reason = source["reason"];
	        this.kind = source["kind"];
	        this.id = source["id"];
	        this.legacy_name = source["legacy_name"];
	        this.message = source["message"];
	    }
	}
	export class PlanItem {
	    source_key: string;
	    source_kind: string;
	    source_id?: string;
	    source_name?: string;
	    source_root?: string;
	    folder?: string;
	    proposed_id?: string;
	    action: string;
	    reason?: string;
	    message?: string;
	    evidence?: string;
	    source_hash: string;
	    document_hash?: string;
	    lore_hashes?: Record<string, string>;
	    media_hash?: string;
	    writes?: string[];
	    references?: ReferenceIssue[];
	
	    static createFrom(source: any = {}) {
	        return new PlanItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_key = source["source_key"];
	        this.source_kind = source["source_kind"];
	        this.source_id = source["source_id"];
	        this.source_name = source["source_name"];
	        this.source_root = source["source_root"];
	        this.folder = source["folder"];
	        this.proposed_id = source["proposed_id"];
	        this.action = source["action"];
	        this.reason = source["reason"];
	        this.message = source["message"];
	        this.evidence = source["evidence"];
	        this.source_hash = source["source_hash"];
	        this.document_hash = source["document_hash"];
	        this.lore_hashes = source["lore_hashes"];
	        this.media_hash = source["media_hash"];
	        this.writes = source["writes"];
	        this.references = this.convertValues(source["references"], ReferenceIssue);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Plan {
	    hash: string;
	    category: string;
	    destination_library_id: string;
	    destination_root?: string;
	    destination_status: string;
	    conversion_version: string;
	    items: PlanItem[];
	    ready: number;
	    adoptable: number;
	    conflicted: number;
	    quarantined: number;
	    unsupported: number;
	    reference_issues: number;
	    writes: number;
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Plan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.category = source["category"];
	        this.destination_library_id = source["destination_library_id"];
	        this.destination_root = source["destination_root"];
	        this.destination_status = source["destination_status"];
	        this.conversion_version = source["conversion_version"];
	        this.items = this.convertValues(source["items"], PlanItem);
	        this.ready = source["ready"];
	        this.adoptable = source["adoptable"];
	        this.conflicted = source["conflicted"];
	        this.quarantined = source["quarantined"];
	        this.unsupported = source["unsupported"];
	        this.reference_issues = source["reference_issues"];
	        this.writes = source["writes"];
	        this.errors = source["errors"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Resolution {
	    source_key: string;
	    destination_id?: string;
	    new_id?: boolean;
	    actor?: string;
	    evidence?: string;
	
	    static createFrom(source: any = {}) {
	        return new Resolution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_key = source["source_key"];
	        this.destination_id = source["destination_id"];
	        this.new_id = source["new_id"];
	        this.actor = source["actor"];
	        this.evidence = source["evidence"];
	    }
	}
	export class Request {
	    category: string;
	    destination_library_id: string;
	    browser_records: any[];
	    legacy_roots: string[];
	    resolutions: Resolution[];
	    operation_id?: string;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.destination_library_id = source["destination_library_id"];
	        this.browser_records = source["browser_records"];
	        this.legacy_roots = source["legacy_roots"];
	        this.resolutions = this.convertValues(source["resolutions"], Resolution);
	        this.operation_id = source["operation_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Result {
	    operation_id: string;
	    status: string;
	    plan_hash: string;
	    plan: Plan;
	    items: ItemResult[];
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation_id = source["operation_id"];
	        this.status = source["status"];
	        this.plan_hash = source["plan_hash"];
	        this.plan = this.convertValues(source["plan"], Plan);
	        this.items = this.convertValues(source["items"], ItemResult);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RollbackItem {
	    item_id: string;
	    state: string;
	    reason?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new RollbackItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item_id = source["item_id"];
	        this.state = source["state"];
	        this.reason = source["reason"];
	        this.message = source["message"];
	    }
	}
	export class RollbackResult {
	    operation_id: string;
	    status: string;
	    items: RollbackItem[];
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new RollbackResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation_id = source["operation_id"];
	        this.status = source["status"];
	        this.items = this.convertValues(source["items"], RollbackItem);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace library {
	
	export class Item {
	    id: string;
	    library_id: string;
	    category: string;
	    name: string;
	    path: string;
	    format: string;
	    source: string;
	    external: boolean;
	    editable: boolean;
	    available: boolean;
	    description?: string;
	    personality?: string;
	    scenario?: string;
	    first_mes?: string;
	    mes_example?: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.library_id = source["library_id"];
	        this.category = source["category"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.format = source["format"];
	        this.source = source["source"];
	        this.external = source["external"];
	        this.editable = source["editable"];
	        this.available = source["available"];
	        this.description = source["description"];
	        this.personality = source["personality"];
	        this.scenario = source["scenario"];
	        this.first_mes = source["first_mes"];
	        this.mes_example = source["mes_example"];
	    }
	}
	export class InventoryRoot {
	    id: string;
	    category: string;
	    path: string;
	    role: number;
	    origin: string;
	    format: string;
	    default_save: boolean;
	    created_at: string;
	    available: boolean;
	    scanned: boolean;
	    stale: boolean;
	    scan_error?: string;
	    scanned_at?: string;
	    items: Item[];
	
	    static createFrom(source: any = {}) {
	        return new InventoryRoot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.path = source["path"];
	        this.role = source["role"];
	        this.origin = source["origin"];
	        this.format = source["format"];
	        this.default_save = source["default_save"];
	        this.created_at = source["created_at"];
	        this.available = source["available"];
	        this.scanned = source["scanned"];
	        this.stale = source["stale"];
	        this.scan_error = source["scan_error"];
	        this.scanned_at = source["scanned_at"];
	        this.items = this.convertValues(source["items"], Item);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class OperationResult {
	    status: string;
	    category?: string;
	    item_id?: string;
	    target?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.category = source["category"];
	        this.item_id = source["item_id"];
	        this.target = source["target"];
	        this.message = source["message"];
	    }
	}
	export class Registration {
	    id: string;
	    category: string;
	    path: string;
	    role: number;
	    origin: string;
	    format: string;
	    default_save: boolean;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Registration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.path = source["path"];
	        this.role = source["role"];
	        this.origin = source["origin"];
	        this.format = source["format"];
	        this.default_save = source["default_save"];
	        this.created_at = source["created_at"];
	    }
	}
	export class View {
	    id: string;
	    category: string;
	    path: string;
	    role: number;
	    origin: string;
	    format: string;
	    default_save: boolean;
	    created_at: string;
	    exists: boolean;
	    readable: boolean;
	    writable: boolean;
	    save_root: boolean;
	    available: boolean;
	
	    static createFrom(source: any = {}) {
	        return new View(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.path = source["path"];
	        this.role = source["role"];
	        this.origin = source["origin"];
	        this.format = source["format"];
	        this.default_save = source["default_save"];
	        this.created_at = source["created_at"];
	        this.exists = source["exists"];
	        this.readable = source["readable"];
	        this.writable = source["writable"];
	        this.save_root = source["save_root"];
	        this.available = source["available"];
	    }
	}

}

export namespace lorebookstore {
	
	export class ImportItem {
	    id: string;
	    action: string;
	    path?: string;
	    detail?: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.action = source["action"];
	        this.path = source["path"];
	        this.detail = source["detail"];
	    }
	}
	export class ImportResult {
	    status: string;
	    cutover: boolean;
	    root_id?: string;
	    created: number;
	    adopted: number;
	    skipped: number;
	    conflicts: number;
	    items: ImportItem[];
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.cutover = source["cutover"];
	        this.root_id = source["root_id"];
	        this.created = source["created"];
	        this.adopted = source["adopted"];
	        this.skipped = source["skipped"];
	        this.conflicts = source["conflicts"];
	        this.items = this.convertValues(source["items"], ImportItem);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Item {
	    id: string;
	    library_id: string;
	    root: string;
	    path: string;
	    origin: string;
	    editable: boolean;
	    available: boolean;
	    updated_at?: string;
	    entry: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.library_id = source["library_id"];
	        this.root = source["root"];
	        this.path = source["path"];
	        this.origin = source["origin"];
	        this.editable = source["editable"];
	        this.available = source["available"];
	        this.updated_at = source["updated_at"];
	        this.entry = source["entry"];
	    }
	}
	export class RootStatus {
	    id: string;
	    path: string;
	    origin: string;
	    available: boolean;
	    writable: boolean;
	    save_root: boolean;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new RootStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.origin = source["origin"];
	        this.available = source["available"];
	        this.writable = source["writable"];
	        this.save_root = source["save_root"];
	        this.count = source["count"];
	    }
	}
	export class ListResult {
	    cutover: boolean;
	    cutover_at?: string;
	    count: number;
	    items: Item[];
	    roots: RootStatus[];
	    conflicts: string[];
	    ignored: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new ListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cutover = source["cutover"];
	        this.cutover_at = source["cutover_at"];
	        this.count = source["count"];
	        this.items = this.convertValues(source["items"], Item);
	        this.roots = this.convertValues(source["roots"], RootStatus);
	        this.conflicts = source["conflicts"];
	        this.ignored = source["ignored"];
	        this.errors = source["errors"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OwnerIdentity {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new OwnerIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class OwnerMigrationItem {
	    id: string;
	    action: string;
	    owner_id?: string;
	    owner_name?: string;
	    detail?: string;
	
	    static createFrom(source: any = {}) {
	        return new OwnerMigrationItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.action = source["action"];
	        this.owner_id = source["owner_id"];
	        this.owner_name = source["owner_name"];
	        this.detail = source["detail"];
	    }
	}
	export class OwnerMigrationResult {
	    status: string;
	    version?: string;
	    migrated: number;
	    already: number;
	    ambiguous: number;
	    unknown: number;
	    failed: number;
	    restored?: number;
	    skipped?: number;
	    items: OwnerMigrationItem[];
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new OwnerMigrationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.version = source["version"];
	        this.migrated = source["migrated"];
	        this.already = source["already"];
	        this.ambiguous = source["ambiguous"];
	        this.unknown = source["unknown"];
	        this.failed = source["failed"];
	        this.restored = source["restored"];
	        this.skipped = source["skipped"];
	        this.items = this.convertValues(source["items"], OwnerMigrationItem);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SaveResult {
	    status: string;
	    id: string;
	    path?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.id = source["id"];
	        this.path = source["path"];
	        this.message = source["message"];
	    }
	}

}

export namespace main {
	
	export class AllocationAdvice {
	    use_gpu: boolean;
	    gpu_layers: number;
	    gpu_memory_mb: number;
	    system_ram_mb: number;
	    reduce_batch: boolean;
	    reduce_context: boolean;
	    unload_models: string[];
	    reason: string;
	    confidence: number;
	
	    static createFrom(source: any = {}) {
	        return new AllocationAdvice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.use_gpu = source["use_gpu"];
	        this.gpu_layers = source["gpu_layers"];
	        this.gpu_memory_mb = source["gpu_memory_mb"];
	        this.system_ram_mb = source["system_ram_mb"];
	        this.reduce_batch = source["reduce_batch"];
	        this.reduce_context = source["reduce_context"];
	        this.unload_models = source["unload_models"];
	        this.reason = source["reason"];
	        this.confidence = source["confidence"];
	    }
	}
	export class RecommendedInfo {
	    sampler: string;
	    scheduler: string;
	    steps: number;
	    cfg: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new RecommendedInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sampler = source["sampler"];
	        this.scheduler = source["scheduler"];
	        this.steps = source["steps"];
	        this.cfg = source["cfg"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class AssetInfo {
	    name: string;
	    path: string;
	    type: string;
	    family: string;
	    format: string;
	    size_mb: number;
	    compatibility: string;
	    files?: Record<string, string>;
	    recommended: RecommendedInfo;
	
	    static createFrom(source: any = {}) {
	        return new AssetInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.type = source["type"];
	        this.family = source["family"];
	        this.format = source["format"];
	        this.size_mb = source["size_mb"];
	        this.compatibility = source["compatibility"];
	        this.files = source["files"];
	        this.recommended = this.convertValues(source["recommended"], RecommendedInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AutoCaptureConfig {
	    enabled: boolean;
	    minConfidence: number;
	    maxPerMessage: number;
	    triggerWords?: string[];
	
	    static createFrom(source: any = {}) {
	        return new AutoCaptureConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.minConfidence = source["minConfidence"];
	        this.maxPerMessage = source["maxPerMessage"];
	        this.triggerWords = source["triggerWords"];
	    }
	}
	export class DependencyItem {
	    role: string;
	    label: string;
	    state: string;
	    path?: string;
	    fileName?: string;
	    url?: string;
	    source?: string;
	    sizeMB?: number;
	    note?: string;
	    mismatch?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new DependencyItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.label = source["label"];
	        this.state = source["state"];
	        this.path = source["path"];
	        this.fileName = source["fileName"];
	        this.url = source["url"];
	        this.source = source["source"];
	        this.sizeMB = source["sizeMB"];
	        this.note = source["note"];
	        this.mismatch = source["mismatch"];
	        this.error = source["error"];
	    }
	}
	export class DependencyReport {
	    model: string;
	    family: string;
	    known: boolean;
	    items: DependencyItem[];
	    missing: string[];
	    downloadDir: string;
	    totalSizeMB: number;
	    canDownload: boolean;
	    offline: boolean;
	    message?: string;
	    consentText: string;
	
	    static createFrom(source: any = {}) {
	        return new DependencyReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.family = source["family"];
	        this.known = source["known"];
	        this.items = this.convertValues(source["items"], DependencyItem);
	        this.missing = source["missing"];
	        this.downloadDir = source["downloadDir"];
	        this.totalSizeMB = source["totalSizeMB"];
	        this.canDownload = source["canDownload"];
	        this.offline = source["offline"];
	        this.message = source["message"];
	        this.consentText = source["consentText"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DirectoryFolder {
	    name: string;
	    path: string;
	    model_files: number;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryFolder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.model_files = source["model_files"];
	    }
	}
	export class DirectoryBrowse {
	    path: string;
	    parent: string;
	    at_root: boolean;
	    home: string;
	    drives: string[];
	    folders: DirectoryFolder[];
	    model_files: number;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryBrowse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.parent = source["parent"];
	        this.at_root = source["at_root"];
	        this.home = source["home"];
	        this.drives = source["drives"];
	        this.folders = this.convertValues(source["folders"], DirectoryFolder);
	        this.model_files = source["model_files"];
	        this.kind = source["kind"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DiscoveryCandidate {
	    path: string;
	    name: string;
	    model_files: number;
	    kind: string;
	    kinds: string[];
	
	    static createFrom(source: any = {}) {
	        return new DiscoveryCandidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.model_files = source["model_files"];
	        this.kind = source["kind"];
	        this.kinds = source["kinds"];
	    }
	}
	export class FilteredAssets {
	    vaes: AssetInfo[];
	    loras: AssetInfo[];
	    controlnets: AssetInfo[];
	
	    static createFrom(source: any = {}) {
	        return new FilteredAssets(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vaes = this.convertValues(source["vaes"], AssetInfo);
	        this.loras = this.convertValues(source["loras"], AssetInfo);
	        this.controlnets = this.convertValues(source["controlnets"], AssetInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageWorkflowInfo {
	    name: string;
	    file_name: string;
	    mode: string;
	    model: string;
	    saved_at: string;
	    has_pose: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImageWorkflowInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.file_name = source["file_name"];
	        this.mode = source["mode"];
	        this.model = source["model"];
	        this.saved_at = source["saved_at"];
	        this.has_pose = source["has_pose"];
	    }
	}
	export class InstalledAssetInfo {
	    name: string;
	    original_path: string;
	    installed_path: string;
	    type: string;
	    family: string;
	    format: string;
	    size_mb: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new InstalledAssetInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.original_path = source["original_path"];
	        this.installed_path = source["installed_path"];
	        this.type = source["type"];
	        this.family = source["family"];
	        this.format = source["format"];
	        this.size_mb = source["size_mb"];
	        this.error = source["error"];
	    }
	}
	export class InstallSummary {
	    registered: InstalledAssetInfo[];
	    already_known: InstalledAssetInfo[];
	    failed: InstalledAssetInfo[];
	    total_count: number;
	    registered_count: number;
	    already_known_count: number;
	    failed_count: number;
	
	    static createFrom(source: any = {}) {
	        return new InstallSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.registered = this.convertValues(source["registered"], InstalledAssetInfo);
	        this.already_known = this.convertValues(source["already_known"], InstalledAssetInfo);
	        this.failed = this.convertValues(source["failed"], InstalledAssetInfo);
	        this.total_count = source["total_count"];
	        this.registered_count = source["registered_count"];
	        this.already_known_count = source["already_known_count"];
	        this.failed_count = source["failed_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class LLMModelDetail {
	    name: string;
	    path: string;
	    format: string;
	    size_mb: number;
	    quant: string;
	    parameter_count: string;
	    layer_count: number;
	    max_context: number;
	    family: string;
	    gpu_mb: number;
	    ram_mb: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new LLMModelDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.format = source["format"];
	        this.size_mb = source["size_mb"];
	        this.quant = source["quant"];
	        this.parameter_count = source["parameter_count"];
	        this.layer_count = source["layer_count"];
	        this.max_context = source["max_context"];
	        this.family = source["family"];
	        this.gpu_mb = source["gpu_mb"];
	        this.ram_mb = source["ram_mb"];
	        this.notes = source["notes"];
	    }
	}
	export class LLMModelInfo {
	    name: string;
	    path: string;
	    format: string;
	    size_mb: number;
	
	    static createFrom(source: any = {}) {
	        return new LLMModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.format = source["format"];
	        this.size_mb = source["size_mb"];
	    }
	}
	export class LoadedModelInfo {
	    type: string;
	    name: string;
	    gpu_memory_mb: number;
	    system_ram_mb: number;
	    priority: number;
	    is_gpu_offload: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LoadedModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.gpu_memory_mb = source["gpu_memory_mb"];
	        this.system_ram_mb = source["system_ram_mb"];
	        this.priority = source["priority"];
	        this.is_gpu_offload = source["is_gpu_offload"];
	    }
	}
	export class MemoryEdit {
	    title?: string;
	    summary?: string;
	    tags?: string[];
	    pinned?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MemoryEdit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.summary = source["summary"];
	        this.tags = source["tags"];
	        this.pinned = source["pinned"];
	    }
	}
	export class MemoryFilter {
	    query?: string;
	    scope?: string;
	    status?: string;
	
	    static createFrom(source: any = {}) {
	        return new MemoryFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.scope = source["scope"];
	        this.status = source["status"];
	    }
	}
	export class MemoryRecord {
	    id: string;
	    createdAt: string;
	    lastTriggeredAt?: string;
	    sourceMessages?: string[];
	    rawText?: string;
	    summary?: string;
	    title: string;
	    type: string;
	    scope: string;
	    scopeId?: string;
	    ownerKind?: string;
	    confidence?: number;
	    status?: string;
	    pinned?: boolean;
	    tags?: string[];
	    author?: string;
	    trigger?: string;
	    scoreBreakdown?: Record<string, any>;
	    resolverVersion?: number;
	    version?: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = source["createdAt"];
	        this.lastTriggeredAt = source["lastTriggeredAt"];
	        this.sourceMessages = source["sourceMessages"];
	        this.rawText = source["rawText"];
	        this.summary = source["summary"];
	        this.title = source["title"];
	        this.type = source["type"];
	        this.scope = source["scope"];
	        this.scopeId = source["scopeId"];
	        this.ownerKind = source["ownerKind"];
	        this.confidence = source["confidence"];
	        this.status = source["status"];
	        this.pinned = source["pinned"];
	        this.tags = source["tags"];
	        this.author = source["author"];
	        this.trigger = source["trigger"];
	        this.scoreBreakdown = source["scoreBreakdown"];
	        this.resolverVersion = source["resolverVersion"];
	        this.version = source["version"];
	    }
	}
	export class ModelDirectory {
	    path: string;
	    kind: string;
	    exists: boolean;
	    models: number;
	
	    static createFrom(source: any = {}) {
	        return new ModelDirectory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.kind = source["kind"];
	        this.exists = source["exists"];
	        this.models = source["models"];
	    }
	}
	export class ModelDiscoveryState {
	    enabled: boolean;
	    suppressed: boolean;
	    prompting: boolean;
	    candidate?: DiscoveryCandidate;
	    deferred: number;
	
	    static createFrom(source: any = {}) {
	        return new ModelDiscoveryState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.suppressed = source["suppressed"];
	        this.prompting = source["prompting"];
	        this.candidate = this.convertValues(source["candidate"], DiscoveryCandidate);
	        this.deferred = source["deferred"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ModelFitResult {
	    kind: string;
	    can_fit: boolean;
	    required_mb: number;
	    free_vram_mb: number;
	    total_vram_mb: number;
	    blocking: LoadedModelInfo[];
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelFitResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.can_fit = source["can_fit"];
	        this.required_mb = source["required_mb"];
	        this.free_vram_mb = source["free_vram_mb"];
	        this.total_vram_mb = source["total_vram_mb"];
	        this.blocking = this.convertValues(source["blocking"], LoadedModelInfo);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ResourceSnapshot {
	    gpu_model: string;
	    gpu_vendor: string;
	    total_vram_mb: number;
	    used_vram_mb: number;
	    free_vram_mb: number;
	    total_ram_mb: number;
	    used_ram_mb: number;
	    free_ram_mb: number;
	    memory_tier: string;
	    os: string;
	
	    static createFrom(source: any = {}) {
	        return new ResourceSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gpu_model = source["gpu_model"];
	        this.gpu_vendor = source["gpu_vendor"];
	        this.total_vram_mb = source["total_vram_mb"];
	        this.used_vram_mb = source["used_vram_mb"];
	        this.free_vram_mb = source["free_vram_mb"];
	        this.total_ram_mb = source["total_ram_mb"];
	        this.used_ram_mb = source["used_ram_mb"];
	        this.free_ram_mb = source["free_ram_mb"];
	        this.memory_tier = source["memory_tier"];
	        this.os = source["os"];
	    }
	}
	export class ScanResult {
	    models: AssetInfo[];
	    vaes: AssetInfo[];
	    loras: AssetInfo[];
	    controlnets: AssetInfo[];
	    embeddings: AssetInfo[];
	    upscalers: AssetInfo[];
	    t5xxls: AssetInfo[];
	    clip_ls: AssetInfo[];
	    clip_gs: AssetInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.models = this.convertValues(source["models"], AssetInfo);
	        this.vaes = this.convertValues(source["vaes"], AssetInfo);
	        this.loras = this.convertValues(source["loras"], AssetInfo);
	        this.controlnets = this.convertValues(source["controlnets"], AssetInfo);
	        this.embeddings = this.convertValues(source["embeddings"], AssetInfo);
	        this.upscalers = this.convertValues(source["upscalers"], AssetInfo);
	        this.t5xxls = this.convertValues(source["t5xxls"], AssetInfo);
	        this.clip_ls = this.convertValues(source["clip_ls"], AssetInfo);
	        this.clip_gs = this.convertValues(source["clip_gs"], AssetInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TTSModelInfo {
	    name: string;
	    path: string;
	    format: string;
	    size_mb: number;
	
	    static createFrom(source: any = {}) {
	        return new TTSModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.format = source["format"];
	        this.size_mb = source["size_mb"];
	    }
	}
	export class TTSVoiceOption {
	    id: string;
	    name: string;
	    language?: string;
	    description?: string;
	    available: boolean;
	    detail?: string;
	
	    static createFrom(source: any = {}) {
	        return new TTSVoiceOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.language = source["language"];
	        this.description = source["description"];
	        this.available = source["available"];
	        this.detail = source["detail"];
	    }
	}
	export class TTSVoiceCatalog {
	    engine: string;
	    inherited: boolean;
	    available: boolean;
	    default_voice?: string;
	    speed_supported: boolean;
	    min_speed: number;
	    max_speed: number;
	    detail?: string;
	    voices: TTSVoiceOption[];
	
	    static createFrom(source: any = {}) {
	        return new TTSVoiceCatalog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engine = source["engine"];
	        this.inherited = source["inherited"];
	        this.available = source["available"];
	        this.default_voice = source["default_voice"];
	        this.speed_supported = source["speed_supported"];
	        this.min_speed = source["min_speed"];
	        this.max_speed = source["max_speed"];
	        this.detail = source["detail"];
	        this.voices = this.convertValues(source["voices"], TTSVoiceOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace memory {
	
	export class IdentityDiagnostics {
	    total: number;
	    stable: number;
	    unresolved: number;
	    unresolvedOwners?: string[];
	
	    static createFrom(source: any = {}) {
	        return new IdentityDiagnostics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.stable = source["stable"];
	        this.unresolved = source["unresolved"];
	        this.unresolvedOwners = source["unresolvedOwners"];
	    }
	}

}

export namespace metadata {
	
	export class ExtractedMetadata {
	    prompt: string;
	    negative_prompt: string;
	    steps: number;
	    sampler: string;
	    cfg_scale: number;
	    seed: number;
	    width: number;
	    height: number;
	    model_hash?: string;
	    model_name?: string;
	    scheduler?: string;
	    raw?: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractedMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.negative_prompt = source["negative_prompt"];
	        this.steps = source["steps"];
	        this.sampler = source["sampler"];
	        this.cfg_scale = source["cfg_scale"];
	        this.seed = source["seed"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.model_hash = source["model_hash"];
	        this.model_name = source["model_name"];
	        this.scheduler = source["scheduler"];
	        this.raw = source["raw"];
	    }
	}

}

export namespace modelintel {
	
	export class Profile {
	    path: string;
	    name: string;
	    format: string;
	    architecture?: string;
	    family?: string;
	    tokenizer?: string;
	    chat_template?: string;
	    fingerprint?: string;
	    capabilities: string[];
	    runtimes: string[];
	    evidence: string;
	    source: string;
	    notes?: string[];
	    available: boolean;
	    size_bytes: number;
	    mod_time_ns: number;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.format = source["format"];
	        this.architecture = source["architecture"];
	        this.family = source["family"];
	        this.tokenizer = source["tokenizer"];
	        this.chat_template = source["chat_template"];
	        this.fingerprint = source["fingerprint"];
	        this.capabilities = source["capabilities"];
	        this.runtimes = source["runtimes"];
	        this.evidence = source["evidence"];
	        this.source = source["source"];
	        this.notes = source["notes"];
	        this.available = source["available"];
	        this.size_bytes = source["size_bytes"];
	        this.mod_time_ns = source["mod_time_ns"];
	        this.updated_at = source["updated_at"];
	    }
	}

}

export namespace rag {
	
	export class Document {
	    id: string;
	    name: string;
	    scope: string;
	    scopeId: string;
	    sourceType: string;
	    sourcePath: string;
	    sizeBytes: number;
	    chunkCount: number;
	    lastUsed: string;
	    createdAt: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.scope = source["scope"];
	        this.scopeId = source["scopeId"];
	        this.sourceType = source["sourceType"];
	        this.sourcePath = source["sourcePath"];
	        this.sizeBytes = source["sizeBytes"];
	        this.chunkCount = source["chunkCount"];
	        this.lastUsed = source["lastUsed"];
	        this.createdAt = source["createdAt"];
	        this.enabled = source["enabled"];
	    }
	}

}

export namespace settings {
	
	export class ImageSettings {
	    engine_config: types.EngineConfig;
	    last_model: string;
	    last_model_opts: types.ModelLoadOptions;
	
	    static createFrom(source: any = {}) {
	        return new ImageSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.engine_config = this.convertValues(source["engine_config"], types.EngineConfig);
	        this.last_model = source["last_model"];
	        this.last_model_opts = this.convertValues(source["last_model_opts"], types.ModelLoadOptions);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace types {
	
	export class ArchitectureSpec {
	    Family: string;
	    Layers: number;
	    DModel: number;
	    FFNSize: number;
	    NumHeads: number;
	    NumKVHeads: number;
	    HeadDim: number;
	    VocabSize: number;
	    ContextLength: number;
	    KVCacheType: string;
	    KVCacheBytesPerElement: number;
	    UseGQA: boolean;
	    UseSlidingWindow: boolean;
	    SlidingWindowSize: number;
	    IsMoE: boolean;
	    ExpertCount: number;
	    ExpertUsedCount: number;
	    RoPEDimension: number;
	    RoPEBase: number;
	    RoPEScale: number;
	    QuantType: string;
	    BitsPerWeight: number;
	    FileSizeBytes: number;
	    ParameterCount: number;
	    HasMetadata: boolean;
	    MissingFields: string[];
	
	    static createFrom(source: any = {}) {
	        return new ArchitectureSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Family = source["Family"];
	        this.Layers = source["Layers"];
	        this.DModel = source["DModel"];
	        this.FFNSize = source["FFNSize"];
	        this.NumHeads = source["NumHeads"];
	        this.NumKVHeads = source["NumKVHeads"];
	        this.HeadDim = source["HeadDim"];
	        this.VocabSize = source["VocabSize"];
	        this.ContextLength = source["ContextLength"];
	        this.KVCacheType = source["KVCacheType"];
	        this.KVCacheBytesPerElement = source["KVCacheBytesPerElement"];
	        this.UseGQA = source["UseGQA"];
	        this.UseSlidingWindow = source["UseSlidingWindow"];
	        this.SlidingWindowSize = source["SlidingWindowSize"];
	        this.IsMoE = source["IsMoE"];
	        this.ExpertCount = source["ExpertCount"];
	        this.ExpertUsedCount = source["ExpertUsedCount"];
	        this.RoPEDimension = source["RoPEDimension"];
	        this.RoPEBase = source["RoPEBase"];
	        this.RoPEScale = source["RoPEScale"];
	        this.QuantType = source["QuantType"];
	        this.BitsPerWeight = source["BitsPerWeight"];
	        this.FileSizeBytes = source["FileSizeBytes"];
	        this.ParameterCount = source["ParameterCount"];
	        this.HasMetadata = source["HasMetadata"];
	        this.MissingFields = source["MissingFields"];
	    }
	}
	export class EngineConfig {
	    binary_path: string;
	    models_dir: string;
	    output_dir: string;
	    temp_dir: string;
	    port: number;
	
	    static createFrom(source: any = {}) {
	        return new EngineConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.binary_path = source["binary_path"];
	        this.models_dir = source["models_dir"];
	        this.output_dir = source["output_dir"];
	        this.temp_dir = source["temp_dir"];
	        this.port = source["port"];
	    }
	}
	export class ImageSource {
	    world_id?: string;
	    scenario_id?: string;
	    character_id?: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.world_id = source["world_id"];
	        this.scenario_id = source["scenario_id"];
	        this.character_id = source["character_id"];
	    }
	}
	export class ImageMetadata {
	    file: string;
	    prompt: string;
	    negative_prompt: string;
	    model: string;
	    family: string;
	    width: number;
	    height: number;
	    steps: number;
	    seed: number;
	    sampler: string;
	    cfg_scale: number;
	    mode: string;
	    created_at: string;
	    source: ImageSource;
	    options?: ImageOptions;
	
	    static createFrom(source: any = {}) {
	        return new ImageMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.prompt = source["prompt"];
	        this.negative_prompt = source["negative_prompt"];
	        this.model = source["model"];
	        this.family = source["family"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.steps = source["steps"];
	        this.seed = source["seed"];
	        this.sampler = source["sampler"];
	        this.cfg_scale = source["cfg_scale"];
	        this.mode = source["mode"];
	        this.created_at = source["created_at"];
	        this.source = this.convertValues(source["source"], ImageSource);
	        this.options = this.convertValues(source["options"], ImageOptions);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageResult {
	    path: string;
	    thumb_path: string;
	    metadata: ImageMetadata;
	    image_data?: string;
	    images?: string[];
	    video_data?: string;
	    video_mime?: string;
	    image_metadatas?: ImageMetadata[];
	    gallery_error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.thumb_path = source["thumb_path"];
	        this.metadata = this.convertValues(source["metadata"], ImageMetadata);
	        this.image_data = source["image_data"];
	        this.images = source["images"];
	        this.video_data = source["video_data"];
	        this.video_mime = source["video_mime"];
	        this.image_metadatas = this.convertValues(source["image_metadatas"], ImageMetadata);
	        this.gallery_error = source["gallery_error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LoRAStrength {
	    name: string;
	    strength: number;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LoRAStrength(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.strength = source["strength"];
	        this.enabled = source["enabled"];
	    }
	}
	export class ImageOptions {
	    prompt: string;
	    negative_prompt: string;
	    width: number;
	    height: number;
	    steps: number;
	    seed: number;
	    sampler: string;
	    scheduler: string;
	    cfg_scale: number;
	    guidance: number;
	    slg_scale: number;
	    skip_layer_start: number;
	    skip_layer_end: number;
	    eta: number;
	    strength: number;
	    init_image_path: string;
	    mask_path: string;
	    model_name: string;
	    vae_path: string;
	    clip_l_path: string;
	    clip_g_path: string;
	    t5xxl_path: string;
	    control_net_path: string;
	    control_image_path: string;
	    control_strength: number;
	    batch_count: number;
	    mode: string;
	    job_id?: string;
	    outpaint_margin?: number;
	    outpaint_direction?: string;
	    upscale_factor?: number;
	    loras?: LoRAStrength[];
	    enable_hr: boolean;
	    hr_scale: number;
	    hr_upscaler: string;
	    denoising_strength: number;
	    resize_mode?: string;
	    style_strength?: number;
	    preserve_colors?: boolean;
	    match_histogram?: boolean;
	    latent_noise?: number;
	    mask_blur?: number;
	    mask_invert?: boolean;
	    mask_feather?: number;
	    masked_content?: string;
	    edge_behavior?: string;
	    upscale_method?: string;
	    upscale_sharpen?: string;
	    upscale_anti_artifact?: string;
	    tiling: boolean;
	    clip_skip: number;
	    video_frames: number;
	    fps: number;
	    video_motion_strength?: number;
	    video_seed_mode?: string;
	    video_denoise?: number;
	    video_loop?: boolean;
	    output_format?: string;
	    end_image_path?: string;
	    cache_mode: string;
	    cache_option: string;
	    prediction?: string;
	    rng?: string;
	    sampler_rng?: string;
	    flow_shift?: number;
	    vae_tile_overlap?: number;
	    disable_metadata?: boolean;
	    cfg_rescale?: number;
	    face_restore?: number;
	    s_churn?: number;
	    s_tmin?: number;
	    s_tmax?: number;
	    s_noise?: number;
	    high_noise_steps?: number;
	    high_noise_cfg_scale?: number;
	    high_noise_guidance?: number;
	    high_noise_slg_scale?: number;
	    high_noise_eta?: number;
	    high_noise_skip_layer_start?: number;
	    high_noise_skip_layer_end?: number;
	    vae_tiling?: boolean;
	    temporal_tiling?: boolean;
	    flash_attention?: boolean;
	    offload_to_cpu?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImageOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.negative_prompt = source["negative_prompt"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.steps = source["steps"];
	        this.seed = source["seed"];
	        this.sampler = source["sampler"];
	        this.scheduler = source["scheduler"];
	        this.cfg_scale = source["cfg_scale"];
	        this.guidance = source["guidance"];
	        this.slg_scale = source["slg_scale"];
	        this.skip_layer_start = source["skip_layer_start"];
	        this.skip_layer_end = source["skip_layer_end"];
	        this.eta = source["eta"];
	        this.strength = source["strength"];
	        this.init_image_path = source["init_image_path"];
	        this.mask_path = source["mask_path"];
	        this.model_name = source["model_name"];
	        this.vae_path = source["vae_path"];
	        this.clip_l_path = source["clip_l_path"];
	        this.clip_g_path = source["clip_g_path"];
	        this.t5xxl_path = source["t5xxl_path"];
	        this.control_net_path = source["control_net_path"];
	        this.control_image_path = source["control_image_path"];
	        this.control_strength = source["control_strength"];
	        this.batch_count = source["batch_count"];
	        this.mode = source["mode"];
	        this.job_id = source["job_id"];
	        this.outpaint_margin = source["outpaint_margin"];
	        this.outpaint_direction = source["outpaint_direction"];
	        this.upscale_factor = source["upscale_factor"];
	        this.loras = this.convertValues(source["loras"], LoRAStrength);
	        this.enable_hr = source["enable_hr"];
	        this.hr_scale = source["hr_scale"];
	        this.hr_upscaler = source["hr_upscaler"];
	        this.denoising_strength = source["denoising_strength"];
	        this.resize_mode = source["resize_mode"];
	        this.style_strength = source["style_strength"];
	        this.preserve_colors = source["preserve_colors"];
	        this.match_histogram = source["match_histogram"];
	        this.latent_noise = source["latent_noise"];
	        this.mask_blur = source["mask_blur"];
	        this.mask_invert = source["mask_invert"];
	        this.mask_feather = source["mask_feather"];
	        this.masked_content = source["masked_content"];
	        this.edge_behavior = source["edge_behavior"];
	        this.upscale_method = source["upscale_method"];
	        this.upscale_sharpen = source["upscale_sharpen"];
	        this.upscale_anti_artifact = source["upscale_anti_artifact"];
	        this.tiling = source["tiling"];
	        this.clip_skip = source["clip_skip"];
	        this.video_frames = source["video_frames"];
	        this.fps = source["fps"];
	        this.video_motion_strength = source["video_motion_strength"];
	        this.video_seed_mode = source["video_seed_mode"];
	        this.video_denoise = source["video_denoise"];
	        this.video_loop = source["video_loop"];
	        this.output_format = source["output_format"];
	        this.end_image_path = source["end_image_path"];
	        this.cache_mode = source["cache_mode"];
	        this.cache_option = source["cache_option"];
	        this.prediction = source["prediction"];
	        this.rng = source["rng"];
	        this.sampler_rng = source["sampler_rng"];
	        this.flow_shift = source["flow_shift"];
	        this.vae_tile_overlap = source["vae_tile_overlap"];
	        this.disable_metadata = source["disable_metadata"];
	        this.cfg_rescale = source["cfg_rescale"];
	        this.face_restore = source["face_restore"];
	        this.s_churn = source["s_churn"];
	        this.s_tmin = source["s_tmin"];
	        this.s_tmax = source["s_tmax"];
	        this.s_noise = source["s_noise"];
	        this.high_noise_steps = source["high_noise_steps"];
	        this.high_noise_cfg_scale = source["high_noise_cfg_scale"];
	        this.high_noise_guidance = source["high_noise_guidance"];
	        this.high_noise_slg_scale = source["high_noise_slg_scale"];
	        this.high_noise_eta = source["high_noise_eta"];
	        this.high_noise_skip_layer_start = source["high_noise_skip_layer_start"];
	        this.high_noise_skip_layer_end = source["high_noise_skip_layer_end"];
	        this.vae_tiling = source["vae_tiling"];
	        this.temporal_tiling = source["temporal_tiling"];
	        this.flash_attention = source["flash_attention"];
	        this.offload_to_cpu = source["offload_to_cpu"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageJob {
	    id: string;
	    status: string;
	    progress: number;
	    step: number;
	    total_steps: number;
	    preview_image?: string;
	    options?: ImageOptions;
	    result?: ImageResult;
	    error?: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.step = source["step"];
	        this.total_steps = source["total_steps"];
	        this.preview_image = source["preview_image"];
	        this.options = this.convertValues(source["options"], ImageOptions);
	        this.result = this.convertValues(source["result"], ImageResult);
	        this.error = source["error"];
	        this.created_at = source["created_at"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class MemoryEstimate {
	    GPUMB: number;
	    RAMMB: number;
	    TotalMB: number;
	    ContextSizeUsed: number;
	    GPULayersUsed: number;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new MemoryEstimate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GPUMB = source["GPUMB"];
	        this.RAMMB = source["RAMMB"];
	        this.TotalMB = source["TotalMB"];
	        this.ContextSizeUsed = source["ContextSizeUsed"];
	        this.GPULayersUsed = source["GPULayersUsed"];
	        this.Notes = source["Notes"];
	    }
	}
	export class ModelInfo {
	    name: string;
	    path: string;
	    family: string;
	    capabilities: string[];
	    resolution_hint: string;
	    files: Record<string, string>;
	    format: string;
	    compatibility: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.family = source["family"];
	        this.capabilities = source["capabilities"];
	        this.resolution_hint = source["resolution_hint"];
	        this.files = source["files"];
	        this.format = source["format"];
	        this.compatibility = source["compatibility"];
	    }
	}
	export class ModelLoadOptions {
	    model: string;
	    model_path?: string;
	    job_id?: string;
	    family?: string;
	    vae: string;
	    clip_l: string;
	    clip_g: string;
	    t5xxl: string;
	    llm: string;
	    diffusion_model: string;
	    backend?: string;
	    params_backend?: string;
	    cache_mode?: string;
	    cache_option?: string;
	    vae_tiling?: boolean;
	    temporal_tiling?: boolean;
	    circular?: boolean;
	    flash_attention?: boolean;
	    offload_to_cpu?: boolean;
	    disable_metadata?: boolean;
	    max_vram?: number;
	    mmap?: boolean;
	    prediction?: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelLoadOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.model_path = source["model_path"];
	        this.job_id = source["job_id"];
	        this.family = source["family"];
	        this.vae = source["vae"];
	        this.clip_l = source["clip_l"];
	        this.clip_g = source["clip_g"];
	        this.t5xxl = source["t5xxl"];
	        this.llm = source["llm"];
	        this.diffusion_model = source["diffusion_model"];
	        this.backend = source["backend"];
	        this.params_backend = source["params_backend"];
	        this.cache_mode = source["cache_mode"];
	        this.cache_option = source["cache_option"];
	        this.vae_tiling = source["vae_tiling"];
	        this.temporal_tiling = source["temporal_tiling"];
	        this.circular = source["circular"];
	        this.flash_attention = source["flash_attention"];
	        this.offload_to_cpu = source["offload_to_cpu"];
	        this.disable_metadata = source["disable_metadata"];
	        this.max_vram = source["max_vram"];
	        this.mmap = source["mmap"];
	        this.prediction = source["prediction"];
	    }
	}
	export class ModelProfile {
	    Path: string;
	    Family: string;
	    Variant: string;
	    MaxContext: number;
	    UsableContext: number;
	    RopeMode: string;
	    RopeFactor: number;
	    RopeBase: number;
	    Template: string;
	    Tokenizer: string;
	    StopSequences: string[];
	    LayerCount: number;
	    ParameterCount: number;
	    FileType: string;
	    FileSizeBytes: number;
	    BitsPerWeight: number;
	    ArchitectureSpec?: ArchitectureSpec;
	    HasVision: boolean;
	    HasAudioInput: boolean;
	    ProjectorKind: string;
	    MMPProjPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Family = source["Family"];
	        this.Variant = source["Variant"];
	        this.MaxContext = source["MaxContext"];
	        this.UsableContext = source["UsableContext"];
	        this.RopeMode = source["RopeMode"];
	        this.RopeFactor = source["RopeFactor"];
	        this.RopeBase = source["RopeBase"];
	        this.Template = source["Template"];
	        this.Tokenizer = source["Tokenizer"];
	        this.StopSequences = source["StopSequences"];
	        this.LayerCount = source["LayerCount"];
	        this.ParameterCount = source["ParameterCount"];
	        this.FileType = source["FileType"];
	        this.FileSizeBytes = source["FileSizeBytes"];
	        this.BitsPerWeight = source["BitsPerWeight"];
	        this.ArchitectureSpec = this.convertValues(source["ArchitectureSpec"], ArchitectureSpec);
	        this.HasVision = source["HasVision"];
	        this.HasAudioInput = source["HasAudioInput"];
	        this.ProjectorKind = source["ProjectorKind"];
	        this.MMPProjPath = source["MMPProjPath"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PromptSettings {
	    PromptTemplate: string;
	    CustomJinjaTemplate: string;
	    SystemPromptOverride: string;
	    UserPrefix: string;
	    AssistantPrefix: string;
	    StopSequences: string[];
	
	    static createFrom(source: any = {}) {
	        return new PromptSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PromptTemplate = source["PromptTemplate"];
	        this.CustomJinjaTemplate = source["CustomJinjaTemplate"];
	        this.SystemPromptOverride = source["SystemPromptOverride"];
	        this.UserPrefix = source["UserPrefix"];
	        this.AssistantPrefix = source["AssistantPrefix"];
	        this.StopSequences = source["StopSequences"];
	    }
	}
	export class RuntimeSettings {
	    Threads: number;
	    BatchSize: number;
	    ContextSize: number;
	    GPULayers: number;
	    RopeMode: string;
	    RopeFactor: number;
	    RopeBase: number;
	    FlashAttention: boolean;
	    TensorSplit: number[];
	    MainGPU: number;
	    OffloadKQV: boolean;
	    UseMMap: boolean;
	    UseMLock: boolean;
	    VocabOnly: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Threads = source["Threads"];
	        this.BatchSize = source["BatchSize"];
	        this.ContextSize = source["ContextSize"];
	        this.GPULayers = source["GPULayers"];
	        this.RopeMode = source["RopeMode"];
	        this.RopeFactor = source["RopeFactor"];
	        this.RopeBase = source["RopeBase"];
	        this.FlashAttention = source["FlashAttention"];
	        this.TensorSplit = source["TensorSplit"];
	        this.MainGPU = source["MainGPU"];
	        this.OffloadKQV = source["OffloadKQV"];
	        this.UseMMap = source["UseMMap"];
	        this.UseMLock = source["UseMLock"];
	        this.VocabOnly = source["VocabOnly"];
	    }
	}

}

