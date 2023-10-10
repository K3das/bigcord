const GrafanaPanel = (props: { dashboard: string, refresh: string, panelId: number }) => {
    return (
        <iframe src={`${import.meta.env.VITE_GRAFANA_BASE_URL ?? "http://localhost:3001"}/d-solo/${props.dashboard}?orgId=1&refresh=${props.refresh}&theme=light&panelId=${props.panelId}`} class={"w-full h-96 border-0"}></iframe>
    )
}

export default GrafanaPanel