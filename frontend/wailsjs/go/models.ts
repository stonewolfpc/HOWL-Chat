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
	    outpaint_margin?: number;
	    outpaint_direction?: string;
	    upscale_factor?: number;
	    loras?: string[];
	    enable_hr: boolean;
	    hr_scale: number;
	    hr_upscaler: string;
	    denoising_strength: number;
	    tiling: boolean;
	    clip_skip: number;
	    video_frames: number;
	    fps: number;
	    cache_mode: string;
	    cache_option: string;
	    prediction?: string;
	    rng?: string;
	    sampler_rng?: string;
	    flow_shift?: number;
	    vae_tile_overlap?: number;
	    disable_metadata?: boolean;
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
	        this.outpaint_margin = source["outpaint_margin"];
	        this.outpaint_direction = source["outpaint_direction"];
	        this.upscale_factor = source["upscale_factor"];
	        this.loras = source["loras"];
	        this.enable_hr = source["enable_hr"];
	        this.hr_scale = source["hr_scale"];
	        this.hr_upscaler = source["hr_upscaler"];
	        this.denoising_strength = source["denoising_strength"];
	        this.tiling = source["tiling"];
	        this.clip_skip = source["clip_skip"];
	        this.video_frames = source["video_frames"];
	        this.fps = source["fps"];
	        this.cache_mode = source["cache_mode"];
	        this.cache_option = source["cache_option"];
	        this.prediction = source["prediction"];
	        this.rng = source["rng"];
	        this.sampler_rng = source["sampler_rng"];
	        this.flow_shift = source["flow_shift"];
	        this.vae_tile_overlap = source["vae_tile_overlap"];
	        this.disable_metadata = source["disable_metadata"];
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
	    vae: string;
	    clip_l: string;
	    clip_g: string;
	    t5xxl: string;
	    diffusion_model: string;
	    backend?: string;
	    params_backend?: string;
	    vae_tiling?: boolean;
	    temporal_tiling?: boolean;
	    circular?: boolean;
	    flash_attention?: boolean;
	    offload_to_cpu?: boolean;
	    disable_metadata?: boolean;
	    prediction?: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelLoadOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.model_path = source["model_path"];
	        this.vae = source["vae"];
	        this.clip_l = source["clip_l"];
	        this.clip_g = source["clip_g"];
	        this.t5xxl = source["t5xxl"];
	        this.diffusion_model = source["diffusion_model"];
	        this.backend = source["backend"];
	        this.params_backend = source["params_backend"];
	        this.vae_tiling = source["vae_tiling"];
	        this.temporal_tiling = source["temporal_tiling"];
	        this.circular = source["circular"];
	        this.flash_attention = source["flash_attention"];
	        this.offload_to_cpu = source["offload_to_cpu"];
	        this.disable_metadata = source["disable_metadata"];
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

