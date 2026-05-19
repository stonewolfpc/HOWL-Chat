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

export namespace types {
	
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
	
	    static createFrom(source: any = {}) {
	        return new ImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.thumb_path = source["thumb_path"];
	        this.metadata = this.convertValues(source["metadata"], ImageMetadata);
	        this.image_data = source["image_data"];
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

