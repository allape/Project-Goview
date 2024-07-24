<script lang="ts">
  let loading: boolean = false;
  let step: number = 0;

  let timerId = -1;

  const LoadingChars = ['-', '\\', '|', '/'];

  export let onClick: () => undefined | Promise<unknown> = () => undefined;

  async function handleClick() {
    if (loading) {
      return;
    }
    loading = true;
    clearInterval(timerId);
    timerId = setInterval(() => {
      step = (step + 1) % LoadingChars.length;
    }, 100);
    try {
      await onClick();
    } finally {
      loading = false;
    }
  }
</script>

<style lang="scss">
  .button {
    position: relative;
    .loadingText {
      width: 100%;
      height: 100%;
      overflow: hidden;
      white-space: nowrap;
      position: absolute;
      top: 0;
      left: 0;
      font-family: monospace;
    }
  }
</style>

<button {...$$props} class:button={true} disabled={loading} on:click={handleClick}>
  <span class="loadingText" style:opacity={loading ? 1 : 0}>
    {LoadingChars[step]}
  </span>
  <span style:opacity={loading ? 0 : 1}>
    <slot></slot>
  </span>
</button>
