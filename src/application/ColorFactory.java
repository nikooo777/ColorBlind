package application;

import java.util.concurrent.ThreadLocalRandom;

import javafx.scene.paint.Color;

/*
 * colors to use:
 * 0;65;1 //Dark green
 * 75;56;0 //Dark brown
 * 141;161;222 //bluish
 * 169;148;224 //pinkish
 * --------
 * 0;174;0 //green
 * 110;0;0 //Dark red
 * ------------------
 * 31;31;248 //Blue
 * 255;34;51 //red
 * 119;123;140 //grey
 * 238;213;30
 *
 */
public class ColorFactory
{
	private static final Color DARKGREEN = new Color(normalize(0), normalize(65), normalize(1), 1);
	private static final Color DARKBROWN = new Color(normalize(75), normalize(56), normalize(0), 1);
	private static final Color BLUISH = new Color(normalize(141), normalize(161), normalize(222), 1);
	private static final Color PINKISH = new Color(normalize(169), normalize(148), normalize(224), 1);
	private static final Color GREEN = new Color(normalize(0), normalize(174), normalize(0), 1);
	private static final Color DARKRED = new Color(normalize(110), normalize(0), normalize(0), 1);
	private static final Color BLUE = new Color(normalize(31), normalize(31), normalize(248), 0.8);
	private static final Color RED = new Color(normalize(255), normalize(34), normalize(51), 1);
	private static final Color GREY = new Color(normalize(119), normalize(123), normalize(140), 1);
	private static final Color ORANGE = new Color(normalize(242), normalize(142), normalize(244), 1);

	private static double normalize(double value)
	{
		return value / 255.;
	}

	public static Color randomNormalColor()
	{
		// final double rand = ThreadLocalRandom.current().nextGaussian();
		final double rand = ThreadLocalRandom.current().nextDouble();
		if (rand < 0.1)
			return PINKISH;
		else if (rand > 0.9)
			return BLUISH;
		else if (rand < 0.2)
			return GREY;
		else if (rand > 0.8)
			return ORANGE;
		else if (rand < 0.3)
			return GREEN;
		else if (rand > 0.7)
			return RED;
		else if (rand < 0.4)
			return DARKRED;
		else if (rand > 0.6)
			return DARKBROWN;
		else if (rand < 0.5)
			return DARKGREEN;
		else// (rand > 0.5)
			return BLUE;
	}

	public static Color randomC()
	{
		return new Color(normalize(255), normalize(0), normalize(157), 1);
	}
}
