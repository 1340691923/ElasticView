import { SidebarLightThemeEnum } from "@/enums/ThemeEnum";

// 辅助函数：将十六进制颜色转换为 RGB
function hexToRgb(hex: string): [number, number, number] {
  const bigint = parseInt(hex.slice(1), 16);
  return [(bigint >> 16) & 255, (bigint >> 8) & 255, bigint & 255];
}

// 辅助函数：将 RGB 转换为十六进制颜色
function rgbToHex(r: number, g: number, b: number): string {
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

// 辅助函数：调整颜色亮度
function adjustBrightness(hex: string, factor: number): string {
  const rgb = hexToRgb(hex);
  const newRgb = rgb.map((val) =>
    Math.max(0, Math.min(255, Math.round(val + (255 - val) * factor)))
  ) as [number, number, number];
  return rgbToHex(...newRgb);
}

export function generateThemeColors(primary: string) {
  const colors: Record<string, string> = {
    primary,
  };

  // 生成浅色变体
  for (let i = 1; i <= 9; i++) {
    const factor = i * 0.1;
    colors[`primary-light-${i}`] = adjustBrightness(primary, factor);
  }

  // 生成深色变体
  colors["primary-dark-2"] = adjustBrightness(primary, -0.2);

  return colors;
}

export function applyTheme(colors: Record<string, string>) {
  const el = document.documentElement;

  Object.entries(colors).forEach(([key, value]) => {
    el.style.setProperty(`--el-color-${key}`, value);
  });
}

/**
 * 切换暗黑模式
 *
 * @param isDark 是否启用暗黑模式
 */
export function toggleDarkMode(isDark: boolean) {
  if (isDark) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
}

/**
 * 切换浅色主题下的侧边栏颜色方案
 *
 * @param isBlue 布尔值，表示是否开启深蓝色侧边栏颜色方案
 */
export function toggleLightModeSidebarColorScheme(isDarkBlueDark: boolean) {
  if (isDarkBlueDark) {
    document.documentElement.classList.add(SidebarLightThemeEnum.DARKBLUE);
  } else {
    document.documentElement.classList.remove(SidebarLightThemeEnum.DARKBLUE);
  }
}
