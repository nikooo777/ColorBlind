package application;

import java.util.Random;

import javafx.application.Application;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.StackPane;
import javafx.stage.Screen;
import javafx.stage.Stage;

public class Main extends Application
{

	static Random rand;

	@Override
	public void start(final Stage primaryStage /* which is the main window */)
	{
		// set the title
		primaryStage.setTitle("Reverse Colorblind Message Encrypter");
		// root layout to be inserted in the scene
		final BorderPane root = new BorderPane();
		final StackPane layout = new StackPane();
		root.setCenter(layout);

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

		// useless button
		final Button button = new Button("Open the 1337 window");
		final ColorBlindWindow cbw = new ColorBlindWindow();
		button.setOnAction(e ->
		{
			cbw.show();
		});

		layout.getChildren().add(button);
		// pop up the window
		primaryStage.show();

		// // define the biggest circle which is the hidden round
		// final Circle biggest = new Circle(primaryScreenBounds.getHeight() / 3);
		// biggest.setCenterX((primaryScreenBounds.getWidth() / 2.) * rand.nextDouble());
		// biggest.setCenterY(primaryScreenBounds.getHeight() / 2.);
		// biggest.setOpacity(0);
		// final Circle biggest = CircleFactory.shapeCircle(scene);

		// generate all circles
		// final List<Circle> circles = CircleFactory.generateCircles(scene);

		// debug message
		// System.out.println("Successfully generated all circles!");

		// change the color to the circles that have to pop
		// drawSecrets(biggest, circles);

		// add all the circles to the layout
		// root.getChildren().addAll(circles);
		// root.getChildren().add(biggest);
	}

	// main method
	public static void main(final String[] args)
	{
		// initialize the random
		rand = new Random();
		rand.setSeed(System.currentTimeMillis());

		// launch javaFX
		launch(args);
	}
}
