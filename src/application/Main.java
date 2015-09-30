package application;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.ThreadLocalRandom;

import javafx.application.Application;
import javafx.scene.Scene;
import javafx.scene.layout.BorderPane;
import javafx.scene.paint.Color;
import javafx.scene.shape.Circle;
import javafx.stage.Stage;

public class Main extends Application
{
	static Random rand;

	@Override
	public void start(Stage primaryStage)
	{
		try
		{
			final BorderPane root = new BorderPane();
			final Scene scene = new Scene(root, 1080, 720);
			scene.getStylesheets().add(getClass().getResource("application.css").toExternalForm());
			primaryStage.setScene(scene);
			final Circle biggest = new Circle(720. / 3);
			biggest.setCenterX(1080 / 2.);
			biggest.setCenterY(720 / 2.);
			biggest.setOpacity(0.1);
			final List<Circle> circles = new ArrayList<>();
			circles.add(circleFactory(scene.getWidth(), scene.getHeight(), biggest));
			for (int i = 0; i < 2000; i++)
			{
				Circle c;
				boolean condition = false;
				do
				{
					c = circleFactory(scene.getWidth(), scene.getHeight(), biggest);
					for (final Circle circle : circles)
					{
						final double distance = distance(c, circle);
						if (circle.getRadius() + c.getRadius() > distance)
						{
							condition = false;
							break;
						}
						condition = true;
					}
				} while (condition == false);
				circles.add(c);
				System.out.println("Added circle! [" + i + "]");
			}
			root.getChildren().addAll(circles);
			root.getChildren().add(biggest);
			root.setStyle("-fx-background-color: #DADADA");
			primaryStage.show();
		} catch (final Exception e)
		{
			e.printStackTrace();
		}
	}

	private double distance(Circle a, final Circle b)
	{
		final double distance = Math.sqrt(Math.pow(a.getCenterX() - b.getCenterX(), 2) + Math.pow(a.getCenterY() - b.getCenterY(), 2));
		return distance;
	}

	public Circle circleFactory(double width, double height, Circle big)
	{

		final Circle c = new Circle();
		double r;
		double x;
		double y;

		r = ThreadLocalRandom.current().nextDouble(5, 20);// rand. * 30;
		x = width * rand.nextDouble();
		y = height * rand.nextDouble();
		c.setRadius(r);
		if (Intersects(c, big))
			c.setFill(Color.RED);
		else
			c.setFill(randomColor());
		c.setCenterX(x);
		c.setCenterY(y);
		c.setOpacity(ThreadLocalRandom.current().nextDouble(0.5, 1));
		return c;
	}

	public boolean Intersects(Circle a, Circle b)
	{

		final double minrad = Math.min(a.getRadius(), b.getRadius());
		final double maxrad = Math.max(a.getRadius(), b.getRadius());
		final double distance = distance(a, b);
		System.out.println("Distance: " + distance + " minrad: " + minrad + " maxrad: " + maxrad);
		if (distance < Math.abs(maxrad + minrad) && distance(a, b) > Math.abs(maxrad - minrad))
			return true;
		return false;

	}

	public Color randomColor()
	{
		return new Color(rand.nextDouble(), rand.nextDouble(), rand.nextDouble(), 1);
	}

	public static void main(String[] args)
	{
		rand = new Random();
		rand.setSeed(System.currentTimeMillis());

		launch(args);
	}
}
