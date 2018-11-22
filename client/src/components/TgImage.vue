<template>
	<div class="tg-image">
		<div :style="`background-image: url(${image.file})`" />
	</div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { Image } from '@/models';

@Component
export default class TgMenu extends Vue {
  @Prop({default: {}})
  public image!: Image;
}
</script>

<style lang="scss">
	.tg-image {
		position: relative;
		width: calc(100% - 10px);
		margin: 5px;
		&:after {
			content: "";
			display: block;
			padding-bottom: 100%;
		}
		> div {
			position: absolute;
			width: 100%;
			height: 100%;
      background-size: cover;
		}
	}
	@media screen and (min-width: 400px) {	.tg-image {	width: calc(100%/2 - 10px);	}	}
	@media screen and (min-width: 600px) {	.tg-image {	width: calc(100%/3 - 10px);	}	}
	@media screen and (min-width: 850px) {	.tg-image {	width: calc(100%/4 - 10px);	}	}
	@media screen and (min-width: 1100px) {	.tg-image {	width: calc(100%/5 - 10px);	}	}
	@media screen and (min-width: 1400px) {	.tg-image {	width: calc(100%/6 - 10px);	}	}
	@media screen and (min-width: 1700px) {	.tg-image {	width: calc(100%/7 - 10px);	}	}
</style>
