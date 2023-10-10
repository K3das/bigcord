export enum State {
    INVALID,
    FRESH,
    IN_PROGRESS,
    CRASHED,
    COMPLETED
}

export enum JobState {
    INVALID,
    RUNNING,
    FAILED,
    CANCELED
}

export enum ChannelType {
    GUILD_TEXT = 0,
    DM = 1,
    GUILD_VOICE = 2,
    GROUP_DM = 3,
    GUILD_CATEGORY = 4,
    GUILD_ANNOUNCEMENT = 5,
    ANNOUNCEMENT_THREAD = 10,
    PUBLIC_THREAD = 11,
    PRIVATE_THREAD = 12,
    GUILD_STAGE_VOICE = 13,
    GUILD_DIRECTORY = 14,
    GUILD_FORUM = 15,
    GUILD_MEDIA = 16,
}

export type Channel = {
    id: string;
    guild_id: string;
    type: ChannelType;
    name?: string;
    parent_id?: string;
    has_access?: boolean;
    state?: State;
    offset?: string;
    job_state?: JobState;
    job_error?: string;
};

export type GetChannelsResponse = {
    channels: Channel[];
    lost_jobs: Map<string, Job>
}

export type Guild = {
    id: string;
    name: string;
    channels: Channel[];
};

export type Job = {
    ID: string;
    State: JobState;
    JobError: string;
};

export type GenericResponse<T> = {
    success: boolean;
    data?: T;
    error?: string;
};

const Endpoints = Object.freeze({
    DISCORD_GUILD: (guildId: string): string => `/discord/guilds/${guildId}`,
    CHANNELS: "/channels/",
    CHANNEL: (channelId: string): string => `/channels/${channelId}`,
});

export class APIClient {
    private readonly baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    private async handleResponse<T>(response: Response): Promise<T> {
        if (!response.ok) {
            throw new Error(`Request failed with status: ${response.status}`);
        }
        
        const responseData: GenericResponse<T> = await response.json();

        if (!responseData.success) {
            throw new Error(responseData.error || "Unknown error");
        }

        return responseData.data;
    }

    async getDiscordGuild(guildId: string): Promise<Guild> {
        const url = this.baseURL + Endpoints.DISCORD_GUILD(guildId);
        const response = await fetch(url);
        return this.handleResponse<Guild>(response);
    }

    async getChannels(): Promise<GetChannelsResponse> {
        const url = this.baseURL + Endpoints.CHANNELS;
        const response = await fetch(url);
        return this.handleResponse<GetChannelsResponse>(response);
    }

    async getChannel(channelId: string): Promise<Channel> {
        const url = this.baseURL + Endpoints.CHANNEL(channelId);
        const response = await fetch(url);
        return this.handleResponse<Channel>(response);
    }
    
    async addChannels(channels: Set<string>): Promise<Channel> {
        const url = this.baseURL + Endpoints.CHANNELS;
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify(Array.from(channels))
        });
        return this.handleResponse<Channel>(response);
    }
}