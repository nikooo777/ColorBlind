package application;

import java.util.concurrent.ThreadLocalRandom;

import javafx.beans.value.ChangeListener;
import javafx.geometry.Pos;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.control.Slider;
import javafx.scene.layout.StackPane;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;
import javafx.scene.shape.Circle;
import javafx.stage.Modality;
import javafx.stage.Stage;

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
	private static final int MAXOPACITY = 1;
	private static final double MINOPACITY = 0.5;
	private static final Color DARKGREEN = new Color(normalize(0), normalize(65), normalize(1), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color DARKBROWN = new Color(normalize(75), normalize(56), normalize(0), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color BLUISH = new Color(normalize(141), normalize(161), normalize(222), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color PINKISH = new Color(normalize(169), normalize(148), normalize(224), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color GREEN = new Color(normalize(0), normalize(174), normalize(0), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color DARKRED = new Color(normalize(110), normalize(0), normalize(0), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color BLUE = new Color(normalize(31), normalize(31), normalize(248), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color RED = new Color(normalize(255), normalize(34), normalize(51), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color GREY = new Color(normalize(119), normalize(123), normalize(140), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color ORANGE = new Color(normalize(242), normalize(142), normalize(244), ThreadLocalRandom.current().nextDouble(MINOPACITY, MAXOPACITY));
	private static final Color MAGENTAMAIN = new Color(normalize(255), normalize(0), normalize(155), 1);
	private static final Color MAGENTASECONDARY = new Color(normalize(255), normalize(0), normalize(170), 1);

	private static boolean toggle = true;

	private static double normalize(final double value)
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

	public static Color RandomDaltonize(final Circle c)
	{
		toggle = !toggle;
		if (toggle)
			return BLUISH;
		else
			return PINKISH;
	}

	/*
	 * -------------------------------------------------------------------
	 * Second method
	 * -------------------------------------------------------------------
	 */
	public static Color magentaMain()
	{
		return MAGENTAMAIN;
	}

	public static Color magentaSecondary()
	{
		return MAGENTASECONDARY;
	}

	public static Color colorPicker()
	{
		final Stage window = new Stage();
		window.initModality(Modality.APPLICATION_MODAL);
		window.setTitle("Color picker");
		window.setMinWidth(250);
		window.setMinHeight(250);

		final Label label = new Label("move the slider as far left as possible so that you can see the inner circle");
		label.setAlignment(Pos.CENTER);

		final Circle outerCircle = new Circle(50);
		final Circle innerCircle = new Circle(30);
		outerCircle.setFill(magentaMain());
		innerCircle.setFill(magentaSecondary());
		// final Shape wholeCircle = Shape.subtract(outerCircle, innerCircle);
		final Slider slider = new Slider(155, 190, 170);
		slider.valueProperty().addListener((ChangeListener<Number>) (ov, old_val, new_val) ->
		{
			innerCircle.setFill(new Color(normalize(255), normalize(0), normalize(new_val.doubleValue()), 1));
		});

		final Button okButton = new Button("Select current");
		okButton.setOnAction(e ->
		{
			window.close();
		});

		final VBox layout = new VBox(20);
		final StackPane stackCircle = new StackPane();
		stackCircle.getChildren().addAll(outerCircle, innerCircle);
		layout.getChildren().addAll(label, stackCircle, slider, okButton);

		final Scene scene = new Scene(layout);
		window.setScene(scene);
		window.show();
		return (Color) innerCircle.getFill();

	}
}
