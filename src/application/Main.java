package application;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.ThreadLocalRandom;

import javafx.application.Application;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.layout.BorderPane;
import javafx.scene.paint.Color;
import javafx.scene.shape.Circle;
import javafx.stage.Screen;
import javafx.stage.Stage;

public class Main extends Application
{
	private static final int MAXOPACITY = 1;
	private static final double MINOPACITY = 0.5;
	private static final int MAXRAD = 20;
	private static final int MINRAD = 5;
	static Random rand;

	@Override
	public void start(Stage primaryStage /* which is the main window */)
	{
		try
		{
			// set the title
			primaryStage.setTitle("Reverse Colorblind Message Encrypter");
			// root layout to be inserted in the scene
			final BorderPane root = new BorderPane();

			// for the sake of fun, get the primary screen size and set the window to be that size without "maximizing" it
			final Rectangle2D primaryScreenBounds = Screen.getPrimary().getVisualBounds();

			// scene containing the root layout
			final Scene scene = new Scene(root, primaryScreenBounds.getWidth(), primaryScreenBounds.getHeight());

			// set Stage boundaries to visible bounds of the main screen
			primaryStage.setX(primaryScreenBounds.getMinX());
			primaryStage.setY(primaryScreenBounds.getMinY());
			primaryStage.setWidth(primaryScreenBounds.getWidth());
			primaryStage.setHeight(primaryScreenBounds.getHeight());

			// apply css styles (such as the background color)
			scene.getStylesheets().add(getClass().getResource("application.css").toExternalForm());

			// set the scene in the stage
			primaryStage.setScene(scene);

			// define the biggest circle which is the hidden round
			final Circle biggest = new Circle(primaryScreenBounds.getHeight() / 3);
			biggest.setCenterX(primaryScreenBounds.getWidth() / 2.);
			biggest.setCenterY(primaryScreenBounds.getHeight() / 2.);
			biggest.setOpacity(0);

			// prepare a list of circles to be filled
			final List<Circle> circles = new ArrayList<>();

			// generate X unique circles
			for (int i = 0; i < 4000; i++)
			{
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
						if (overlap(c, circle))
						{
							condition = false;
							break;
						}
					}
					// keep on generating the next circle until it doesn't overlap with any of the previous ones
				} while (condition == false);

				// at this point the circle is valid so it's added to the list
				circles.add(c);

				// debug message
				System.out.println("Added circle! [" + i + "]");
			}

			// debug message
			System.out.println("Successfully generated all circles!");

			// change the color to the circles that have to pop
			drawSecrets(biggest, circles);

			// add all the circles to the layout
			root.getChildren().addAll(circles);
			root.getChildren().add(biggest);

			// draw the whole scene
			primaryStage.show();

		} catch (final Exception e)
		{
			e.printStackTrace();
		}
	}

	// draws the circles intersecting the main figure (at this stage it's a circle)
	private void drawSecrets(final Circle biggest, final List<Circle> circles)
	{
		for (final Circle c : circles)
		{
			if (Intersects(c, biggest))
				c.setFill(Color.GREEN); // hardcoded color value
		}
	}

	// are the circles overlapping?
	private boolean overlap(Circle a, final Circle b)
	{
		return b.getRadius() + a.getRadius() > distance(a, b);
	}

	// returns the distance from one circle to the other
	private double distance(Circle a, final Circle b)
	{
		final double centerXsquared = Math.pow(a.getCenterX() - b.getCenterX(), 2);
		final double centerYsquared = Math.pow(a.getCenterY() - b.getCenterY(), 2);
		return Math.sqrt(centerXsquared + centerYsquared);
	}

	// random circles factory
	public Circle circleFactory(double width, double height)
	{

		final Circle c = new Circle();
		double radius, x_position, y_position;

		// calculates a radius value between 5 and 20 (hardcoded)
		radius = ThreadLocalRandom.current().nextDouble(MINRAD, MAXRAD);

		// calculates a random position in the monitor
		x_position = width * rand.nextDouble();
		y_position = height * rand.nextDouble();

		// set the properties
		c.setRadius(radius);
		c.setCenterX(x_position);
		c.setCenterY(y_position);

		// set the color of the circle
		c.setFill(randomColor());

		// Finally return the circle
		return c;
	}

	// does the circle intersect another circle?
	public boolean Intersects(Circle a, Circle b)
	{
		// store the min and max radious of the two circles
		final double minrad = Math.min(a.getRadius(), b.getRadius());
		final double maxrad = Math.max(a.getRadius(), b.getRadius());

		// store the distance between the two circles
		final double distance = distance(a, b);

		// before doing an intense calculation quickly check if they're far apart
		if (minrad + maxrad < distance)
			return false;

		// do intense calculations
		if (distance < Math.abs(maxrad + minrad) && distance(a, b) > Math.abs(maxrad - minrad))
			return true;

		// if nothing is found (they are inside the figure) then return false
		return false;
	}

	// returns a random color
	public Color randomColor()
	{
		// RGB values
		final double R = rand.nextDouble();
		final double G = rand.nextDouble();
		final double B = rand.nextDouble();

		// Opacity with 0 being transparent and 1 fully visible
		final double O = ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY);

		return new Color(R, G, B, O);
	}

	// main method
	public static void main(String[] args)
	{
		// initialize the random
		rand = new Random();
		rand.setSeed(System.currentTimeMillis());

		// launch javaFX
		launch(args);
	}
}
