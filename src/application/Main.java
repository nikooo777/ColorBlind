package application;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.ThreadLocalRandom;

import javafx.application.Application;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.layout.BorderPane;
import javafx.scene.shape.Circle;
import javafx.stage.Screen;
import javafx.stage.Stage;

public class Main extends Application
{
	private static final int GENERATEDCIRCLES = 5000;
	private static final int MAXRAD = 15;
	private static final int MINRAD = 5;
	static Random rand;

	@Override
	public void start(Stage primaryStage /* which is the main window */)
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

		// pop up the window
		primaryStage.show();

		// define the biggest circle which is the hidden round
		final Circle biggest = new Circle(primaryScreenBounds.getHeight() / 3);
		biggest.setCenterX((primaryScreenBounds.getWidth() / 2.) * rand.nextDouble());
		biggest.setCenterY(primaryScreenBounds.getHeight() / 2.);
		biggest.setOpacity(0);

		// generate all circles
		final long pasttime = System.currentTimeMillis();
		final List<Circle> circles = generateCircles(scene);
		System.out.println("Generating " + GENERATEDCIRCLES + " took exactly " + (System.currentTimeMillis() - pasttime) / 1000. + " seconds");

		// debug message
		System.out.println("Successfully generated all circles!");

		// change the color to the circles that have to pop
		drawSecrets(biggest, circles);

		// add all the circles to the layout
		root.getChildren().addAll(circles);
		root.getChildren().add(biggest);
	}

	private List<Circle> generateCircles(final Scene scene)
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

	// draws the circles intersecting the main figure (at this stage it's a circle)
	private void drawSecrets(final Circle biggest, final List<Circle> circles)
	{
		for (final Circle c : circles)
		{
			if (CircleMath.Intersects(c, biggest))
			{
				c.setFill(ColorFactory.MagentaSecondary());
				// c.setFill(new Color(0, 0.25, 0, 1));
				// c.setFill(new Color(0, 0, 0.3, 1));
				// c.setFill(Color.GREEN); // hardcoded color value
			}
		}
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
		c.setFill(ColorFactory.MagentaMain());

		// Finally return the circle
		return c;
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
