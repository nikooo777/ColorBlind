package application;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ThreadLocalRandom;

import javafx.scene.Scene;
import javafx.scene.shape.Circle;

public class CircleFactory
{
	private static final int GENERATEDCIRCLES = 4000;
	private static final int MAXRAD = 15;
	private static final int MINRAD = 5;

	public static List<Circle> generateCircles(final Scene scene)
	{
		long curtime;
		int flagcount = 0;
		// prepare a list of circles to be filled
		final List<Circle> circles = new ArrayList<>();

		// generate X unique circles
		for (int i = 0; i < GENERATEDCIRCLES; i++)
		{
			curtime = System.currentTimeMillis();
			// circle holder
			Circle c;
			// Condition for which the circle is valid
			boolean condition = true;

			// if the circle overlaps to other circles then recreate it
			do
			{
				// generate random circle
				c = circleFactory(scene.getWidth(), scene.getHeight());

				// set the condition to true again
				condition = true;

				// check against all previous circles if there are previous circles
				if (circles.isEmpty())
				{
					break;
				}

				for (final Circle circle : circles)
				{
					if (CircleMath.overlap(c, circle))
					{
						condition = false;
						break;
					}
				}
				// keep on generating the next circle until it doesn't overlap with any of the previous ones
			} while (condition == false);

			// at this point the circle is valid so it's added to the list
			circles.add(c);

			long diff_time;
			if ((diff_time = System.currentTimeMillis() - curtime) > 100)
				flagcount++;

			// debug message
			System.err.println("Added circle! [" + i + "] in " + diff_time + " ms");

			if (flagcount > 10)
				break;
		}
		return circles;
	}

	// random circles factory
	private static Circle circleFactory(double width, double height)
	{

		final Circle c = new Circle();
		double radius, x_position, y_position;

		// calculates a radius value between 5 and 20 (hardcoded)
		radius = ThreadLocalRandom.current().nextDouble(MINRAD, MAXRAD);

		// calculates a random position in the monitor
		x_position = width * Main.rand.nextDouble();
		y_position = height * Main.rand.nextDouble();

		// set the properties
		c.setRadius(radius);
		c.setCenterX(x_position);
		c.setCenterY(y_position);

		// set the color of the circle
		c.setFill(ColorFactory.MagentaMain());

		// Finally return the circle
		return c;
	}
}
